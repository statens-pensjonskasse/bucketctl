package webhook

import (
	"bucketctl/pkg"
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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

	var desiredWebhooks map[string]*ProjectWebhooks
	if err := pkg.ReadConfigFile(file, &desiredWebhooks); err != nil {
		return err
	}

	projectKeys := make([]string, 0, len(desiredWebhooks))
	for p := range desiredWebhooks {
		projectKeys = append(projectKeys, p)
	}
	sort.Strings(projectKeys)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(desiredWebhooks)).WithRemoveWhenDone(true).WithWriter(os.Stderr).Start()
	for _, projectKey := range projectKeys {
		progressBar.Title = projectKey
		actualWebhooks, err := getProjectWebhooks(baseUrl, projectKey, limit, token, true)
		if err != nil {
			return err
		}

		toCreate, toUpdate, toDelete := findWebhooksToChange(desiredWebhooks[projectKey].Webhooks, actualWebhooks.Webhooks)
		for _, w := range toCreate {
			if err := createProjectWebhook(baseUrl, projectKey, token, w); err != nil {
				return err
			}
		}
		for _, w := range toUpdate {
			if err := updateProjectWebhook(baseUrl, projectKey, token, w); err != nil {
				return err
			}
		}
		for _, w := range toDelete {
			if err := deleteProjectWebhook(baseUrl, projectKey, token, w); err != nil {
				return err
			}
		}
		allProjectRepositories, err := repository.GetProjectRepositories(baseUrl, projectKey, limit)
		if err != nil {
			return err
		}
		if desiredWebhooks[projectKey].Repositories == nil {
			desiredWebhooks[projectKey].Repositories = make(map[string]*RepositoryWebhooks)
		}
		for repoSlug := range allProjectRepositories {
			if actualWebhooks.Repositories[repoSlug] == nil {
				actualWebhooks.Repositories[repoSlug] = new(RepositoryWebhooks)
			}
			if desiredWebhooks[projectKey].Repositories[repoSlug] == nil {
				desiredWebhooks[projectKey].Repositories[repoSlug] = new(RepositoryWebhooks)
			}
			toCreate, toUpdate, toDelete := findWebhooksToChange(desiredWebhooks[projectKey].Repositories[repoSlug].Webhooks, actualWebhooks.Repositories[repoSlug].Webhooks)
			for _, w := range toCreate {
				if err := createRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
					return err
				}
			}
			for _, w := range toUpdate {
				if err := updateRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
					return err
				}
			}
			for _, w := range toDelete {
				if err := deleteRepositoryWebhook(baseUrl, projectKey, repoSlug, token, w); err != nil {
					return err
				}
			}

		}
		progressBar.Increment()
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
func sortWebhooksByBestAvailableCandidate(similar []*similarWebhooks) {
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

func createProjectWebhook(baseUrl string, projectKey string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks", baseUrl, projectKey)
	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	if _, err := pkg.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	return nil
}

func updateProjectWebhook(baseUrl string, projectKey string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d", baseUrl, projectKey, webhook.Id)
	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	if _, err := pkg.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	return nil
}

func deleteProjectWebhook(baseUrl string, projectKey string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/webhooks/%d", baseUrl, projectKey, webhook.Id)
	if _, err := pkg.DeleteRequest(url, token, nil); err != nil {
		return err
	}
	return nil
}

func createRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks", baseUrl, projectKey, repoSlug)
	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	if _, err := pkg.PostRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	return nil
}

func updateRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d", baseUrl, projectKey, repoSlug, webhook.Id)
	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	if _, err := pkg.PutRequest(url, token, bytes.NewReader(payload), nil); err != nil {
		return err
	}
	return nil
}

func deleteRepositoryWebhook(baseUrl string, projectKey string, repoSlug string, token string, webhook *types.Webhook) error {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s/repos/%s/webhooks/%d", baseUrl, projectKey, repoSlug, webhook.Id)
	if _, err := pkg.DeleteRequest(url, token, nil); err != nil {
		return err
	}
	return nil
}
