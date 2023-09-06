package git

import (
	"bucketctl/pkg/cmd/repository"
	"bucketctl/pkg/common"
	"bucketctl/pkg/types"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var cloneCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(types.ProjectKeyFlag)
		viper.BindPFlag(types.ProjectKeyFlag, cmd.Flags().Lookup(types.ProjectKeyFlag))
		viper.BindPFlag(types.RepoSlugFlag, cmd.Flags().Lookup(types.RepoSlugFlag))

		viper.BindPFlag(types.IncludeArchivedFlag, cmd.Flags().Lookup(types.IncludeArchivedFlag))
		viper.BindPFlag(types.UpdateFlag, cmd.Flags().Lookup(types.UpdateFlag))
		viper.BindPFlag(types.ForceFlag, cmd.Flags().Lookup(types.ForceFlag))
	},
	Use:   "clone [flags] [path]",
	Short: "Clone repository from origin or all repositories in a project",
	Long:  "",
	Args:  cobra.MinimumNArgs(0),
	RunE:  clone,
}

func init() {
	cloneCmd.Flags().BoolP(types.IncludeArchivedFlag, types.IncludeArchivedFlagShorthand, false, "Include archived repositories")
	cloneCmd.Flags().BoolP(types.UpdateFlag, types.UpdateFlagShorthand, false, "Switch to default branch and pull from origin")
	cloneCmd.Flags().BoolP(types.ForceFlag, types.ForceFlagShorthand, false, "Force sync of default branch")
}

func clone(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(types.BaseUrlFlag)
	gitUrl := viper.GetString(types.GitUrlFlag)
	projectKey := viper.GetString(types.ProjectKeyFlag)
	repoSlug := viper.GetString(types.RepoSlugFlag)
	token := viper.GetString(types.TokenFlag)
	limit := viper.GetInt(types.LimitFlag)

	includeArchived := viper.GetBool(types.IncludeArchivedFlag)
	update := viper.GetBool(types.UpdateFlag)
	force := viper.GetBool(types.ForceFlag)

	var basePath string
	if len(args) >= 1 {
		basePath = args[0]
	} else if repoSlug == "" {
		basePath = projectKey
	} else {
		basePath = "."
	}

	repos := make(map[string]*repository.Repository)
	if repoSlug == "" {
		projectRepos, err := repository.GetProjectRepositories(baseUrl, projectKey, token, limit)
		cobra.CheckErr(err)
		repos = projectRepos
	} else {
		repoInfo, err := repository.GetRepository(baseUrl, projectKey, repoSlug, token)
		cobra.CheckErr(err)
		repos[repoSlug] = repoInfo
	}

	sortedRepoSlugs := common.GetLexicallySortedKeys(repos)
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(sortedRepoSlugs)).WithRemoveWhenDone(true).Start()
	for _, slug := range sortedRepoSlugs {
		progressBar.UpdateTitle(slug)
		if !includeArchived && repos[slug].Archived {
			progressBar.Increment()
			continue
		}
		repoPath := filepath.Join(basePath, slug)
		url := fmt.Sprintf("%s/%s/%s.git", gitUrl, projectKey, slug)

		home, _ := os.UserHomeDir()

		// Try to clone without Auth
		var err error
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL: url,
		})

		if err != nil {
			if strings.Contains(err.Error(), "ssh: handshake failed") {
				sshFile := home + "/.ssh/id_ed25519"
				pterm.Info.Println(sshFile)

				auth, err := getSSHPublicKeys(sshFile)
				pterm.Warning.Println(err)
				if err != nil {
					return err
				}
				_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
					URL:  url,
					Auth: auth,
				})
			}

			if errors.Is(err, git.ErrRepositoryAlreadyExists) {
				if !update {
					pterm.Info.Println("🚀 Skipping already existing repository " + projectKey + "/" + slug)
				} else {
					defaultBranch, err := repository.GetDefaultBranch(baseUrl, projectKey, slug, token)
					if err != nil {
						return errors.New("Error fetching default branch for " + projectKey + "/" + slug)
					}
					if err := syncRefWithRemote(repoPath, defaultBranch.Id, force); err == nil {
						pterm.Info.Println("🔝 Synced " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId + " with origin")
					} else if errors.Is(err, git.NoErrAlreadyUpToDate) {
						pterm.Info.Println("👍 Branch " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId + " already up to date")
					} else {
						pterm.Warning.Println(err.Error() + ": " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId)
					}
				}
			} else if errors.Is(err, transport.ErrEmptyRemoteRepository) {
				pterm.Warning.Println(err.Error() + ": " + projectKey + "/" + slug)
			} else {
				pterm.Error.Println(err)
			}
		} else {
			pterm.Info.Println("⭐️ Cloned " + projectKey + "/" + slug)
		}
		progressBar.Increment()
	}
	return nil
}
