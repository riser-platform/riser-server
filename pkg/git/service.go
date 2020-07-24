package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

var (
	// ErrNoChanges indicates that there were no changes to commit
	ErrNoChanges       = errors.New("no changes to commit")
	KubeSSHMountPath   = "/etc/riser/kube/ssh/identity"
	KubeSSHTargetPath  = "/etc/riser/ssh/identity"
	KubeSSHKeyFileMode = os.FileMode(0400)
)

const (
	// TODO: After adding auth consider making this dynamic e.g. "riser-server (initiated by johndoe@acme.org)"
	commitName  = "riser-server"
	commitEmail = "riser-server@tempuri.org"
	remoteName  = "origin"

	// TODO: Consider making configurable - the main scenario is large repos that take a long time for the initial clone
	gitExecTimeoutSeconds = 30 * time.Second
)

type RepoSettings struct {
	URL         string
	Branch      string
	LocalGitDir string
}

type Repo interface {
	Commit(message string, files []core.ResourceFile) error
	Push() error
	ResetHardRemote() error
	// Lock locks the repo. Be sure to call Unlock when your work is completed.
	Lock()
	// Unlock unlocks the repo.
	Unlock()
}

type repo struct {
	settings *RepoSettings
	sync.Mutex
}

// NewRepo creates a new reference to a repo. There should only be one instance running per git folder.
// WARNING: any pending changes or local commits will be lost
func NewRepo(repoSettings RepoSettings) (Repo, error) {
	repo := &repo{
		settings: &repoSettings,
		Mutex:    sync.Mutex{},
	}

	// Terrible hack due to https://github.com/kubernetes/kubernetes/issues/57923
	// Essentially Kubernetes mounts secrets with too open permissions for SSH, regardless of permission settings, if the container is
	// running as a lower privalaged user.
	_, err := os.Stat(KubeSSHMountPath)
	if err == nil {
		_, err = execWithContext(context.Background(), exec.Command("sh", "-c", fmt.Sprintf("cp %s %s", KubeSSHMountPath, KubeSSHTargetPath)))
		if err != nil {
			return nil, errors.Wrap(err, "Error copying SSH key")
		}
		err := os.Chmod(KubeSSHTargetPath, KubeSSHKeyFileMode)
		if err != nil {
			return nil, errors.Wrap(err, "Error setting permissions on SSH key")
		}
	}

	err = repo.init()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *repo) Commit(message string, files []core.ResourceFile) error {
	err := processFiles(repo.settings.LocalGitDir, files)
	if err != nil {
		return err
	}
	err = repo.addAll()
	if err != nil {
		return err
	}
	err = repo.execGitCmd("commit", "-m", message, "--author", fmt.Sprintf("%s <%s>", commitName, commitEmail))
	if err != nil && isNoChangesErr(err) {
		return ErrNoChanges
	}
	return err
}

func (repo *repo) addAll() error {
	return repo.execGitCmd("add", "--all")
}

func (repo *repo) Push() error {
	return repo.execGitCmd("push")
}

func (repo *repo) init() error {
	err := os.RemoveAll(repo.settings.LocalGitDir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error cleaning git dir: %s", repo.settings.LocalGitDir))
	}

	err = util.EnsureDir(util.EnsureTrailingSlash(repo.settings.LocalGitDir), workspaceFilePerm)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error ensuring git dir: %s", repo.settings.LocalGitDir))
	}

	return repo.clone()
}

func (repo *repo) clone() error {
	return repo.execGitCmd("clone", "--branch", repo.settings.Branch, "--single-branch", "--depth=1", repo.settings.URL, repo.settings.LocalGitDir)
}

func (repo *repo) fetch() error {
	return repo.execGitCmd("fetch", "-f", remoteName, repo.settings.Branch)
}

// ResetHardRemote ensures that the remote is up-to-date. Pending commits will be lost.
func (repo *repo) ResetHardRemote() error {
	// Always fetch before resetting to the remote to ensure that we're up-to-date
	err := repo.fetch()
	if err != nil {
		return err
	}

	return repo.execGitCmd("reset", "--hard", fmt.Sprintf("%s/%s", remoteName, repo.settings.Branch))
}

func (repo *repo) execGitCmd(args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), gitExecTimeoutSeconds)
	defer cancel()
	cmd := repo.buildGitCmd(ctx, args...)
	_, err := execWithContext(ctx, cmd)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("git %s", args))
	}

	return nil
}

func (repo *repo) buildGitCmd(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repo.settings.LocalGitDir
	return cmd
}
