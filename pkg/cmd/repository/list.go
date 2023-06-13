package repository

import (
	"bucketctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
	"strconv"
)

var (
	key string
)

func prettyFormatRepositories(reposMap map[string]*Repository) [][]string {
	repos := make([]string, 0, len(reposMap))
	for s := range reposMap {
		repos = append(repos, s)
	}
	sort.Strings(repos)

	var data [][]string
	data = append(data, []string{"ID", "Slug", "State", "Public", "Archived"})
	for _, slug := range repos {
		row := []string{strconv.Itoa(reposMap[slug].Id), slug, reposMap[slug].StatusMessage, strconv.FormatBool(reposMap[slug].Public), strconv.FormatBool(reposMap[slug].Archived)}
		data = append(data, row)
	}

	return data
}

func init() {
	listRepositoriesCmd.Flags().StringVarP(&key, "key", "k", "", "Project key")
	listRepositoriesCmd.MarkFlagRequired("key")
}

var listRepositoriesCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
	},
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List repositories in a given project",
	RunE:    listRepositories,
}

func listRepositories(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var limit = viper.GetInt("limit")

	repos, err := GetProjectRepositories(baseUrl, projectKey, limit)
	if err != nil {
		return err
	}

	return pkg.PrintData(repos, prettyFormatRepositories)
}
