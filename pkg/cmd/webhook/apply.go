package webhook

import (
	"bucketctl/pkg"
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
)

var (
	fileName string
)

var applyWebhooksCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("file", cmd.Flags().Lookup("file"))
		viper.BindPFlag("include-repos", cmd.Flags().Lookup("include-repos"))
	},
	Use:  "apply",
	RunE: applyWebhooks,
}

func init() {
	applyWebhooksCmd.Flags().StringVarP(&fileName, "file", "f", "", "Webhooks file")
	applyWebhooksCmd.Flags().Bool("include-repos", false, "Include repositories")

	applyWebhooksCmd.MarkFlagRequired("file")
}

func applyWebhooks(cmd *cobra.Command, args []string) error {
	file := viper.GetString("file")
	baseUrl := viper.GetString("baseUrl")
	limit := viper.GetInt("limit")
	token := viper.GetString("token")
	includeRepos := viper.GetBool("include-repos")

	var desiredWebhooks map[string]*ProjectWebhooks
	if err := pkg.ReadConfigFile(file, &desiredWebhooks); err != nil {
		return err
	}

	for projectKey, desiredProjectWebhooks := range desiredWebhooks {

		actualWebhooks, err := getProjectWebhooks(baseUrl, projectKey, limit, token, includeRepos)
		if err != nil {
			return err
		}

		for repoSlug, desiredRepoWebhooks := range desiredProjectWebhooks.Repositories {
			var webhooksToCreate []*types.Webhook
			var webhooksToDelete []*types.Webhook
			var webhooksToUpdate []*types.Webhook
			actualRepoWebhooks, exists := actualWebhooks.Repositories[repoSlug]
			if !exists {
				// Dersom repoet ikke har noen webhooks må alle lages
				webhooksToCreate = desiredRepoWebhooks.Webhooks
			} else {
				availableWebhooks := make([]*types.Webhook, 0, len(actualRepoWebhooks.Webhooks))
				copy(availableWebhooks, actualRepoWebhooks.Webhooks)

				for _, desiredRepoWebhook := range desiredRepoWebhooks.Webhooks {
					if len(availableWebhooks) == 0 {
						// Hvis vi har gått tom for tilgjengelige webhooks å oppdatere må vi lage den
						webhooksToCreate = append(webhooksToCreate, desiredRepoWebhook)
					} else {
						// Finn webhooken mest lik den ønskede
						sort.SliceStable(availableWebhooks, func(i, j int) bool {
							return desiredRepoWebhook.Similarity(availableWebhooks[i]) > desiredRepoWebhook.Similarity(availableWebhooks[j])
						})
						desiredRepoWebhook.Id = availableWebhooks[0].Id
						webhooksToUpdate = append(webhooksToUpdate, desiredRepoWebhook)
						// Fjern webhooken fra tilgjengelige
						availableWebhooks = availableWebhooks[1:]
					}
				}
				// Resterende webhooks sletter vi
				for _, w := range availableWebhooks {
					webhooksToDelete = append(webhooksToDelete, w)
				}
			}
			pterm.Info.Println(webhooksToCreate)
			pterm.Info.Println(webhooksToDelete)

		}

		pterm.Info.Println(desiredProjectWebhooks)
		pterm.Info.Println(actualWebhooks)
	}

	return nil
}

func findWebhooksToChange(desiredWebhooks []*types.Webhook, actualWebhooks []*types.Webhook) (toCreate []*types.Webhook, toUpdate []*types.Webhook, toDelete []*types.Webhook) {
	var webhooksToCreate []*types.Webhook
	var webhooksToDelete []*types.Webhook
	var webhooksToUpdate []*types.Webhook
	if actualWebhooks == nil || len(actualWebhooks) == 0 {
		// Vi har ingen aktuelle webhooks å ta av og må opprette alle nye
		webhooksToCreate = desiredWebhooks
		return webhooksToCreate, webhooksToUpdate, webhooksToDelete
	}
	availableWebhooks := make([]**types.Webhook, len(actualWebhooks))
	for i := range actualWebhooks {
		availableWebhooks = append(availableWebhooks, &actualWebhooks[i])
	}
	// Finn webhookene som ligner mest på hverandre
	ratedWebhooks := rateCandidateWebhooksSimilarity(desiredWebhooks, availableWebhooks)
	// Begynn med den ønskede webhooken som er mest lik en av de aktuelle
	for r := range ratedWebhooks {
		// Plukker ut de ønskede webhookene som ikke har blitt brukt enda
		sorted := ratedWebhooks[r:]
		// Sorter de resterende ønskede webhookene etter beste tilgjengelige kandidat
		sortWebhooksByAvailableCandidatesSimilarity(sorted)
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
					webhooksToUpdate = append(webhooksToUpdate, updatedWebhook)
				}
				// Sett den brukte webhooken som utilgjengelig
				*candidate.webhook = nil
				// Vi har funnet en kandidat og kan slutte å lete
				break
			} else if i >= len(hasBestCandidate.candidates)-1 {
				// Vi har brukt opp alle kandidater
				webhooksToCreate = append(webhooksToCreate, desired)
			}
		}
	}
	for _, w := range availableWebhooks {
		if w != nil && *w != nil {
			// Dersom vi har noen ubrukte tilgjengelige webhooks så skal disse slettes
			webhooksToDelete = append(webhooksToDelete, *w)
		}
	}

	return webhooksToCreate, webhooksToUpdate, webhooksToDelete
}

type similarWebhooks struct {
	webhook    *types.Webhook
	candidates []*similarCandidates
}

type similarCandidates struct {
	similarity float64
	// Bruker dobbeltpeker for å kunne endre på referansen av den første pekningen uten å måtte loope gjennom alle
	// kandidater for å sette de som utilgjengelige
	webhook **types.Webhook
}

// Kalkulerer likheten av base med alle kandidater
func rateCandidateWebhooksSimilarity(baseWebhooks []*types.Webhook, candidateWebhooks []**types.Webhook) []*similarWebhooks {
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
func sortWebhooksByAvailableCandidatesSimilarity(similar []*similarWebhooks) {
	sort.Slice(similar, func(i, j int) bool {
		var ci, cj int
		// Finner indeksene til de beste tilgjengelige kandidaten til sammenligningsgrunnlag
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
