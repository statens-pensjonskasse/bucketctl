package git

import (
	"bucketctl/pkg/types"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

var (
	key  string
	repo string
)

var Cmd = &cobra.Command{
	Use:   "git",
	Short: "Git commands",
}

func init() {
	Cmd.PersistentFlags().StringVarP(&key, types.ProjectKeyFlag, types.ProjectKeyFlagShorthand, "", "Project key")
	Cmd.PersistentFlags().StringVarP(&repo, types.RepoSlugFlag, types.RepoSlugFlagShorthand, "", "Repository slug")

	Cmd.AddCommand(cloneCmd)
}

func syncRefWithRemote(repoPath string, ref string, force bool) error {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	worktree, err := gitRepo.Worktree()
	if err != nil {
		return err
	}
	err = gitRepo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      force,
	})
	if err != nil {
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(ref),
		Force:  force,
	})
	status, _ := worktree.Status()
	if !status.IsClean() && force {
		head, _ := gitRepo.Head()
		worktree.Reset(&git.ResetOptions{Commit: head.Hash(), Mode: git.HardReset})
	}
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      force,
	})
	return err
}

func getSSHPublicKeys(privateKeyFile string) (publicKeys *ssh.PublicKeys, err error) {
	if _, err := os.Stat(privateKeyFile); err != nil {
		return nil, err
	}

	publicKeys, err = ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if strings.Contains(err.Error(), "empty password") {
		password, err := getKeyPassword(privateKeyFile)
		if err != nil {
			return nil, err
		}
		return ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	}
	return publicKeys, err
}

func getKeyPassword(privateKeyFile string) (string, error) {
	fmt.Printf("Enter password for SSH-key %s: ", privateKeyFile)
	password, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}

	return string(password), nil
}
