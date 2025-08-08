package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket"
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"git.spk.no/infra/bucketctl/pkg/common"
	"git.spk.no/infra/bucketctl/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cloneCmd = &cobra.Command{
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(common.ProjectKeyFlag)
		viper.BindPFlag(common.ProjectKeyFlag, cmd.Flags().Lookup(common.ProjectKeyFlag))
		viper.BindPFlag(common.RepoSlugFlag, cmd.Flags().Lookup(common.RepoSlugFlag))

		viper.BindPFlag(common.IncludeArchivedFlag, cmd.Flags().Lookup(common.IncludeArchivedFlag))
		viper.BindPFlag(common.UpdateFlag, cmd.Flags().Lookup(common.UpdateFlag))
		viper.BindPFlag(common.ForceFlag, cmd.Flags().Lookup(common.ForceFlag))
	},
	Use:   "clone [flags] [path]",
	Short: "Clone repository from origin or all repositories in a permission",
	Long:  "",
	Args:  cobra.MinimumNArgs(0),
	RunE:  clone,
}

func init() {
	cloneCmd.Flags().BoolP(common.IncludeArchivedFlag, common.IncludeArchivedFlagShorthand, false, "Include archived repositories")
	cloneCmd.Flags().BoolP(common.UpdateFlag, common.UpdateFlagShorthand, false, "Switch to default branch and pull from origin")
	cloneCmd.Flags().BoolP(common.ForceFlag, common.ForceFlagShorthand, false, "Force sync of default branch")
}

func clone(cmd *cobra.Command, args []string) error {
	baseUrl := viper.GetString(common.BaseUrlFlag)
	gitUrl := viper.GetString(common.GitUrlFlag)
	projectKey := viper.GetString(common.ProjectKeyFlag)
	repoSlug := viper.GetString(common.RepoSlugFlag)
	token := viper.GetString(common.TokenFlag)
	limit := viper.GetInt(common.LimitFlag)

	includeArchived := viper.GetBool(common.IncludeArchivedFlag)
	update := viper.GetBool(common.UpdateFlag)
	force := viper.GetBool(common.ForceFlag)

	var basePath string
	if len(args) >= 1 {
		basePath = args[0]
	} else if repoSlug == "" {
		basePath = projectKey
	} else {
		basePath = "."
	}

	repos := make(map[string]*types.Repository)
	if repoSlug == "" {
		projectRepos, err := bitbucket.GetProjectRepositoriesMap(baseUrl, projectKey, limit, token)
		cobra.CheckErr(err)
		repos = projectRepos
	} else {
		repoInfo, err := bitbucket.GetRepository(baseUrl, projectKey, repoSlug, token)
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

		// Try to clone without Auth
		var err error
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL: url,
		})

		if err != nil {
			if strings.Contains(err.Error(), "ssh: handshake failed") {
				home, _ := os.UserHomeDir()
				sshFile := home + "/.ssh/id_ed25519"

				auth, err := getSSHPublicKeys(sshFile)
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
					logger.Log("🚀 Skipping already existing repository " + projectKey + "/" + slug)
				} else {
					defaultBranch, err := bitbucket.GetDefaultBranch(baseUrl, projectKey, slug, token)
					if err != nil {
						logger.Err("⚠️ Error fetching default branch for " + projectKey + "/" + slug)
						progressBar.Increment()
						continue
					}
					if err := syncRefWithRemote(repoPath, defaultBranch.Id, force); err == nil {
						logger.Log("🔝 Synced " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId + " with origin")
					} else if errors.Is(err, git.NoErrAlreadyUpToDate) {
						logger.Log("👍 Branch " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId + " already up to date")
					} else {
						logger.Warn(err.Error() + ": " + projectKey + "/" + slug + "/" + defaultBranch.DisplayId)
					}
				}
			} else if errors.Is(err, transport.ErrEmptyRemoteRepository) {
				logger.Warn(err.Error() + ": " + projectKey + "/" + slug)
			} else {
				logger.Err(err.Error())
			}
		} else {
			logger.Log("⭐️ Cloned " + projectKey + "/" + slug)
		}
		progressBar.Increment()
	}
	return nil
}
