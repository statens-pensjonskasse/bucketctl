package v1alpha1

import (
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"reflect"
	"sort"
)

type Webhook struct {
	Id                      int         `json:"-" yaml:"-"`
	Name                    string      `json:"name" yaml:"name"`
	CreatedDate             int         `json:"-" yaml:"-"`
	UpdatedDate             int         `json:"-" yaml:"-"`
	Events                  []string    `json:"events" yaml:"events"`
	Configuration           interface{} `json:"configuration" yaml:"configuration"`
	Url                     string      `json:"url" yaml:"url"`
	Active                  bool        `json:"active" yaml:"active"`
	ScopeType               string      `json:"scopeType" yaml:"scopeType"`
	SslVerificationRequired bool        `yaml:"sslVerificationRequired" yaml:"sslVerificationRequired" yaml:"sslVerificationRequired"`
}

type Webhooks []*Webhook

func FromBitbucketWebhook(bitbucketWebhook *types.Webhook) *Webhook {
	sort.Strings(bitbucketWebhook.Events)
	return &Webhook{
		Id:                      bitbucketWebhook.Id,
		Name:                    bitbucketWebhook.Name,
		CreatedDate:             bitbucketWebhook.CreatedDate,
		UpdatedDate:             bitbucketWebhook.CreatedDate,
		Events:                  bitbucketWebhook.Events,
		Configuration:           bitbucketWebhook.Configuration,
		Url:                     bitbucketWebhook.Url,
		Active:                  bitbucketWebhook.Active,
		ScopeType:               bitbucketWebhook.ScopeType,
		SslVerificationRequired: bitbucketWebhook.SslVerificationRequired,
	}
}

func ToBitbucketWebhook(webhook *Webhook) *types.Webhook {
	return &types.Webhook{
		Id:                      webhook.Id,
		Name:                    webhook.Name,
		CreatedDate:             webhook.CreatedDate,
		UpdatedDate:             webhook.UpdatedDate,
		Events:                  webhook.Events,
		Configuration:           webhook.Configuration,
		Url:                     webhook.Url,
		Active:                  webhook.Active,
		ScopeType:               webhook.ScopeType,
		SslVerificationRequired: webhook.SslVerificationRequired,
	}
}

func FindWebhooksToChange(desired *Webhooks, actual *Webhooks) (toCreate *Webhooks, toUpdate *Webhooks, toDelete *Webhooks) {
	if desired == nil {
		desired = new(Webhooks)
	}
	if actual == nil {
		actual = new(Webhooks)
	}

	toCreate = new(Webhooks)
	toUpdate = new(Webhooks)
	toDelete = new(Webhooks)
	if len(*actual) == 0 {
		// Vi har ingen aktuelle webhooks å ta av og må opprette alle nye
		toCreate = desired
		return toCreate, toUpdate, toDelete
	}
	availableWebhooks := make([]**Webhook, len(*actual))
	for i := range *actual {
		availableWebhooks = append(availableWebhooks, &(*actual)[i])
	}
	// Finn webhookene som ligner mest på hverandre
	ratedWebhooks := rateCandidateWebhooksSimilarity(*desired, availableWebhooks)
	// Begynn med den ønskede webhooken som er mest lik en av de aktuelle
	for r := range ratedWebhooks {
		// Plukker ut de ønskede webhookene som ikke har blitt brukt enda
		sorted := ratedWebhooks[r:]
		// Sorter de resterende ønskede webhookene etter beste tilgjengelige kandidat
		sortWebhooksByBestAvailableCandidate(sorted)
		// Plukker ut den ønskede webhooken med den beste kandidaten
		hasBestCandidate := sorted[0]
		desired := hasBestCandidate.webhook

		for i, candidate := range hasBestCandidate.candidates {
			// Bruk den første (mest like) ledige kandidaten
			if *candidate.webhook != nil {
				if !(*desired).Equivalent(*candidate.webhook) {
					// Dersom kandidaten ikke er ekvivalent med den ønskede må den oppdateres
					updatedWebhook := (*desired).Copy()
					updatedWebhook.Id = (*candidate.webhook).Id
					*toUpdate = append(*toUpdate, updatedWebhook)
				}
				// Sett den brukte webhooken som utilgjengelig
				*candidate.webhook = nil
				// Vi har funnet en kandidat og kan slutte å lete
				break
			} else if i >= len(hasBestCandidate.candidates)-1 {
				// Vi har brukt opp alle kandidater
				*toCreate = append(*toCreate, desired)
			}
		}
	}
	for _, w := range availableWebhooks {
		if w != nil && *w != nil {
			// Dersom vi har noen ubrukte tilgjengelige webhooks så skal disse slettes
			*toDelete = append(*toDelete, *w)
		}
	}

	return toCreate, toUpdate, toDelete
}

type similarWebhooks struct {
	webhook    *Webhook
	candidates []*similarCandidates
}

type similarCandidates struct {
	similarity float64
	// Bruker dobbeltpeker for å kunne endre på referansen av den første pekningen uten å måtte loope gjennom alle
	// kandidater for å sette de som utilgjengelige
	webhook **Webhook
}

// Kalkulerer likheten av base med alle kandidater
func rateCandidateWebhooksSimilarity(baseWebhooks []*Webhook, candidateWebhooks []**Webhook) []*similarWebhooks {
	var similar []*similarWebhooks

	for _, webhook := range baseWebhooks {
		var candidates []*similarCandidates
		for _, candidate := range candidateWebhooks {
			if candidate != nil {
				candidates = append(candidates,
					&similarCandidates{
						similarity: webhook.Similarity(*candidate),
						webhook:    candidate,
					},
				)
			}
		}
		// Sorterer kandidater etter likhet
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].similarity > candidates[j].similarity
		})
		similar = append(similar, &similarWebhooks{webhook: webhook, candidates: candidates})
	}
	return similar
}

// Sorterer rangerte webhooks etter den beste tilgjengelige kandidaten
func sortWebhooksByBestAvailableCandidate(similar []*similarWebhooks) {
	sort.Slice(similar, func(i, j int) bool {
		var ci, cj int
		// Finner indeksene til de beste tilgjengelige kandidatene for sammenligningsgrunnlag
		for wi, w := range similar[i].candidates {
			if *w.webhook != nil {
				ci = wi
				break
			}
		}
		for wj, w := range similar[j].candidates {
			if *w.webhook != nil {
				cj = wj
				break
			}
		}
		return similar[i].candidates[ci].similarity > similar[j].candidates[cj].similarity
	})
}

func (whs *Webhooks) Copy() *Webhooks {
	whsCopy := new(Webhooks)
	for _, wh := range *whs {
		*whsCopy = append(*whsCopy, wh.Copy())
	}
	return whsCopy
}

func (wh *Webhook) Copy() *Webhook {
	copied := &Webhook{
		Id:                      wh.Id,
		Name:                    wh.Name,
		CreatedDate:             wh.CreatedDate,
		UpdatedDate:             wh.UpdatedDate,
		Configuration:           wh.Configuration,
		Url:                     wh.Url,
		Active:                  wh.Active,
		ScopeType:               wh.ScopeType,
		SslVerificationRequired: wh.SslVerificationRequired,
	}
	copied.Events = append(copied.Events, wh.Events...)

	return copied
}

func (whs *Webhooks) toMap() map[string]*Webhook {
	asMap := make(map[string]*Webhook, len(*whs))
	for _, wh := range *whs {
		asMap[wh.Name] = wh
	}
	return asMap
}

func (whs *Webhooks) Equals(cmp *Webhooks) bool {
	if whs == cmp {
		return true
	}
	if cmp == nil {
		return false
	}
	if len(*whs) != len(*cmp) {
		return false
	}
	whsMap := whs.toMap()
	cmpMap := cmp.toMap()
	for name, wh := range whsMap {
		if !wh.Equals(cmpMap[name]) {
			return false
		}
	}
	return true
}

func (wh *Webhook) Equals(cmp *Webhook) bool {
	if !wh.Equivalent(cmp) {
		return false
	}
	if wh.Id != cmp.Id {
		return false
	}
	if wh.CreatedDate != cmp.CreatedDate {
		return false
	}
	if wh.UpdatedDate != cmp.UpdatedDate {
		return false
	}
	return true
}

func (wh *Webhook) Equivalent(cmp *Webhook) bool {
	if wh == cmp {
		return true
	}

	if cmp == nil {
		return false
	}
	if wh.Name != cmp.Name {
		return false
	}
	if wh.Url != cmp.Url {
		return false
	}
	if wh.Active != cmp.Active {
		return false
	}
	if wh.SslVerificationRequired != cmp.SslVerificationRequired {
		return false
	}
	if !reflect.DeepEqual(wh.Configuration, cmp.Configuration) {
		return false
	}
	if len(wh.Events) != len(cmp.Events) {
		return false
	}
	elements := make(map[string]struct{}, len(wh.Events))
	// Create a map with all the (unique) elements of list A as keys
	for _, v := range wh.Events {
		elements[v] = struct{}{}
	}
	// Check that all the elements of list B has a key in the map
	for _, v := range cmp.Events {
		if _, exists := elements[v]; !exists {
			return false
		}
	}
	return true
}

// Similarity Finner gir en score på hvor like to webhooks er mellom 0.0 og 1.0
// Dersom ID er lik antas webhookene å være de samme
func (wh *Webhook) Similarity(candidate *Webhook) float64 {
	if candidate == nil {
		return 0.0
	}
	if wh.Id == candidate.Id {
		return 1.0
	}
	similarityScore := 0.0
	if wh.Name == candidate.Name {
		similarityScore += 0.3
	}
	if wh.Url == candidate.Url {
		similarityScore += 0.1
	}
	if wh.Active == candidate.Active {
		similarityScore += 0.1
	}
	if wh.ScopeType == candidate.ScopeType {
		similarityScore += 0.1
	}
	if wh.SslVerificationRequired == candidate.SslVerificationRequired {
		similarityScore += 0.1
	}
	if reflect.DeepEqual(wh.Configuration, candidate.Configuration) {
		similarityScore += 0.1
	}
	if len(wh.Events) == len(candidate.Events) {
		similarityScore += 0.1
	}
	return similarityScore
}
