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
	pairedWebhooks := pairMostSimilarWebhooks(desiredWebhooks, availableWebhooks)
	// Begynn med den ønskede webhooken som er mest lik en av de aktuelle
	for _, pair := range pairedWebhooks {
		desired := pair.webhook
		pterm.Info.Println((*desired).Name)
		for _, p := range pairedWebhooks {
			pterm.Warning.Println((*p.webhook).Name)
			for _, c := range p.candidates {
				if c.webhook != nil && *c.webhook != nil {
					pterm.Error.Println((*c.webhook).Name)
				}
			}
		}
		for i, candidate := range pair.candidates {
			// Bruk den første (mest like) tilgjegelige aktuelle webhooken
			if candidate.webhook != nil && *candidate.webhook != nil {
				(*desired).Id = (*candidate.webhook).Id
				webhooksToUpdate = append(webhooksToUpdate, *desired)

				// Sett den brukte webhooken som utilgjengelig
				*candidate.webhook = nil
				break
			}
			// Vi har brukt opp alle kandidater
			if i == len(pair.candidates)-1 {
				webhooksToCreate = append(webhooksToCreate, *desired)
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
	webhook    **types.Webhook
	candidates []*similarCandidates
}

type similarCandidates struct {
	similarity float64
	webhook    **types.Webhook
}

// Pairs most similar webhooks. Different bases can have the same candidate
func pairMostSimilarWebhooks(baseWebhooks []*types.Webhook, candidateWebhooks []**types.Webhook) []*similarWebhooks {
	var similar []*similarWebhooks

	for i, webhook := range baseWebhooks {
		candidates := []*similarCandidates{{similarity: -1.0, webhook: nil}}
		// Finner grad av likhet med alle kandidater
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
		similar = append(similar, &similarWebhooks{webhook: &baseWebhooks[i], candidates: candidates})
	}
	// Sorterer base webhooks basert på likhet av beste kandidat
	sort.Slice(similar, func(i, j int) bool {
		return similar[i].candidates[0].similarity > similar[j].candidates[0].similarity
	})
	return similar
}
