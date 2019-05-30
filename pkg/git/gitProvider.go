package git

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/pkg/errors"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	// TODO: After adding auth consider making this dynamic e.g. "riser-server (initiated by johndoe@acme.org)"
	commitName        = "riser-server"
	commitEmail       = "riser-server@tempuri.org"
	remoteName        = "origin"
	workspaceFilePerm = 0755
)

var (
	// ErrNoChanges indicates that there were no changes to commit
	ErrNoChanges = errors.New("no changes to commit")
)

type GitRepoProvider interface {
	Commit(message string, files []core.ResourceFile) (string, error)
	Push() error
	ResetHardRemote() error
}

type RepoSettings struct {
	URL         string
	Branch      string
	Username    string
	Password    string
	LocalGitDir string
}

type gitRepoProvider struct {
	repo         *git.Repository
	repoSettings RepoSettings
}

func NewGitRepoProvider(repoSettings RepoSettings) (GitRepoProvider, error) {
	return handleGetRepository(repoSettings, getRepository)
}

// ResetHardRemote resets the current branch to the remote
func (provider gitRepoProvider) ResetHardRemote() error {
	fetchRefSpec := fmt.Sprintf("refs/heads/%s:refs/remotes/%s/%s", provider.repoSettings.Branch, remoteName, provider.repoSettings.Branch)
	err := provider.repo.Fetch(&git.FetchOptions{
		Auth:       createAuth(provider.repoSettings),
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fetchRefSpec),
		},
		Force: true,
	})

	//	plumbing.ErrObjectNotFound is a work around for https://github.com/src-d/go-git/issues/1151
	if err != nil && err != git.NoErrAlreadyUpToDate && err != plumbing.ErrObjectNotFound {
		return errors.Wrap(err, "error fetching changes from repo")
	}

	remoteRefName := fmt.Sprintf("refs/remotes/%s/%s", remoteName, provider.repoSettings.Branch)
	remoteRef, err := provider.repo.Reference(plumbing.ReferenceName(remoteRefName), true)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error getting remote reference %s", remoteRefName))
	}

	wt, err := provider.repo.Worktree()

	if err != nil {
		return errors.Wrap(err, "error getting worktree")
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(provider.repoSettings.Branch),
		Force:  true,
	})

	if err != nil {
		return errors.Wrap(err, "Error checking out branch")
	}

	err = wt.Reset(&git.ResetOptions{
		Commit: remoteRef.Hash(),
		Mode:   git.HardReset,
	})

	if err != nil {
		return errors.Wrap(err, "Error resetting to remote head")
	}

	return nil
}

func (provider gitRepoProvider) Commit(commitMessage string, files []core.ResourceFile) (string, error) {
	worktree, err := provider.repo.Worktree()
	if err != nil {
		return "", errors.Wrap(err, "Error retrieving worktree")
	}

	for _, file := range files {
		fullFileName := filepath.Join(provider.repoSettings.LocalGitDir, file.Name)
		err = util.EnsureDir(fullFileName, workspaceFilePerm)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(fullFileName, file.Contents, workspaceFilePerm)
		if err != nil {
			return "", errors.Wrap(err, "Error writing file")
		}

		_, err = worktree.Add(file.Name)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("Error adding file \"%s\" to worktree.", file.Name))
		}
	}

	status, err := worktree.Status()
	if err != nil {
		return "", errors.Wrap(err, "Error getting worktree status")
	}

	if status.IsClean() {
		return "", ErrNoChanges
	}

	commitHash, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitName,
			Email: commitEmail,
			When:  time.Now().UTC(),
		},
	})

	if err != nil {
		return "", errors.Wrap(err, "Error creating commit")
	}

	return commitHash.String(), nil
}

func (provider gitRepoProvider) Push() error {
	return provider.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		Auth:       createAuth(provider.repoSettings),
	})
}

func handleGetRepository(repoSettings RepoSettings, getRepoFunc func(RepoSettings) (*git.Repository, error)) (*gitRepoProvider, error) {
	repo, err := getRepoFunc(repoSettings)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, errors.New("Unknown error getting repository")
	}
	provider := &gitRepoProvider{
		repo:         repo,
		repoSettings: repoSettings,
	}
	return provider, nil
}

func getRepository(repoSettings RepoSettings) (*git.Repository, error) {
	var repo *git.Repository
	repo, err := git.PlainOpen(repoSettings.LocalGitDir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			repo, err = git.PlainClone(repoSettings.LocalGitDir, false, createCloneOptions(repoSettings))
			if err != nil {
				return nil, err
			}
			return repo, nil
		}
		return nil, err
	}

	return repo, nil
}

func createCloneOptions(repoSettings RepoSettings) *git.CloneOptions {
	// WARNING: Do not set Depth:1 until https://github.com/src-d/go-git/issues/1151 is fixed
	cloneOptions := &git.CloneOptions{
		URL:               repoSettings.URL,
		ReferenceName:     plumbing.NewBranchReferenceName(repoSettings.Branch),
		RemoteName:        remoteName,
		Auth:              createAuth(repoSettings),
		SingleBranch:      true,
		RecurseSubmodules: git.NoRecurseSubmodules,
	}

	return cloneOptions
}

func createAuth(repoSettings RepoSettings) *http.BasicAuth {
	if len(repoSettings.Username) > 0 {
		return &http.BasicAuth{
			Username: repoSettings.Username,
			Password: repoSettings.Password,
		}
	}
	return nil
}
