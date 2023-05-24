package repository

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gobit/pkg"
	"strconv"
)

var (
	key string
)

func prettyFormatRepositories(repos []Repository) [][]string {
	var data [][]string

	data = append(data, []string{"ID", "Name", "State", "Public", "Archived"})

	for _, repo := range repos {
		row := []string{strconv.Itoa(repo.Id), repo.Name, repo.StatusMessage, strconv.FormatBool(repo.Public), strconv.FormatBool(repo.Archived)}
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
	Short:   "List Bitbucket repositories in a given project",
	RunE:    listRepositories,
}

func listRepositories(cmd *cobra.Command, args []string) error {
	var baseUrl = viper.GetString("baseUrl")
	var projectKey = viper.GetString("key")
	var limit = viper.GetInt("limit")

	repos, err := getRepositories(baseUrl, projectKey, limit)
	if err != nil {
		return err
	}

	pkg.PrintData(repos, prettyFormatRepositories)
	return nil
}
