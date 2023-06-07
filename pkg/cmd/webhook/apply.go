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
	availableWebhooks := make([]*types.Webhook, len(actualWebhooks))
	copy(availableWebhooks, actualWebhooks)

	for _, desiredWebhook := range desiredWebhooks {
		if len(availableWebhooks) == 0 {
			// Vi har gått tom for webhooks å oppdatere og må opprette nye
			webhooksToCreate = append(webhooksToCreate, desiredWebhook)
		} else {
			// Finn den aktuelle webhooken mest lik den ønskede
			mostSimilarIndex := 0
			for i, webhook := range availableWebhooks {
				if desiredWebhook.Similarity(webhook) > desiredWebhook.Similarity(availableWebhooks[mostSimilarIndex]) {
					mostSimilarIndex = i
				}
			}
			// Bruk IDen til den mest like webhooken
			desiredWebhook.Id = availableWebhooks[mostSimilarIndex].Id
			webhooksToUpdate = append(webhooksToUpdate, desiredWebhook)

			// Fjern den beste webhooken fra tilgjengelige
			availableWebhooks[mostSimilarIndex] = availableWebhooks[len(availableWebhooks)-1]
			availableWebhooks = availableWebhooks[:len(availableWebhooks)-1]
		}
	}
	// De resterende tilgjengelige webhookene skal slettes
	for _, w := range availableWebhooks {
		webhooksToDelete = append(webhooksToDelete, w)
	}

	return webhooksToCreate, webhooksToUpdate, webhooksToDelete
}

type similarWebhooks struct {
	similarity float64
	base       *types.Webhook
	candidate  *types.Webhook
}

func pairMostSimilar(base []*types.Webhook, candidate []*types.Webhook) []*similarWebhooks {
	var similar []*similarWebhooks
	for _, b := range base {
		mostSimilar := &similarWebhooks{similarity: 0, base: b}
		for _, c := range candidate {
			// Finn kandidaten mest lik basen
			similarity := b.Similarity(c)
			if similarity > mostSimilar.similarity {
				mostSimilar.similarity = similarity
				mostSimilar.candidate = c
			}
		}
		similar = append(similar, mostSimilar)
	}
	// sorter etter likhet
	sort.Slice(similar, func(i, j int) bool {
		return similar[i].similarity > similar[j].similarity
	})
	return similar
}
