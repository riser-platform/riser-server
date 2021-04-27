package git

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	URL    string
	Branch string
	// BaseWorkspaceDir is the root folder for this repo's workspace. A random folder will be created here.
	BaseWorkspaceDir string
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
	settings     *RepoSettings
	workspaceDir string
	sync.Mutex
}

// InitRepoWorkspace clones a repo reference into the specified folder and returns a new reference to the repo.
// If the desired branch does not exist it will be created and pushed to the remote
// WARNING: Running this against the same instance will result in losing unsynchronized changes
func InitRepoWorkspace(repoSettings RepoSettings) (Repo, error) {
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
	err := processFiles(repo.workspaceDir, files)
	if err != nil {
		return err
	}
	err = repo.addAll()
	if err != nil {
		return err
	}
	_, err = repo.execGitCmd("commit", "-m", message, "--author", fmt.Sprintf("%s <%s>", commitName, commitEmail))
	if err != nil && isNoChangesErr(err) {
		return ErrNoChanges
	}
	return err
}

func (repo *repo) Push() error {
	_, err := repo.execGitCmd("push")
	return err
}

// ResetHardRemote ensures that the remote is up-to-date. Pending commits will be lost.
func (repo *repo) ResetHardRemote() error {
	// Always fetch before resetting to the remote to ensure that we're up-to-date
	err := repo.fetch()
	if err != nil {
		return err
	}

	_, err = repo.execGitCmd("reset", "--hard", fmt.Sprintf("%s/%s", remoteName, repo.settings.Branch))
	return err
}

func (repo *repo) addAll() error {

	_, err := repo.execGitCmd("add", "--all")
	return err
}

func (r *repo) init() error {
	err := util.EnsureDir(util.EnsureTrailingSlash(r.settings.BaseWorkspaceDir), workspaceFilePerm)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error ensuring git dir: %s", r.workspaceDir))
	}

	workspaceDir, err := ioutil.TempDir(r.settings.BaseWorkspaceDir, "repo-*")
	if err != nil {
		return errors.Wrap(err, "Error creating workspace dir")
	}
	r.workspaceDir = workspaceDir

	// Create the branch on demand if needed.
	branchExists, err := r.branchExists()
	if err != nil {
		return err
	}

	if branchExists {
		return r.shallowClone()
	} else {
		err = r.shallowCloneDefaultBranch()
		if err != nil {
			return err
		}
		err = r.createEmptyBranch()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error creating branch %q", r.settings.Branch))
		}

		_, err = r.execGitCmd("push", "--set-upstream", remoteName, r.settings.Branch)
		if err != nil {
			return err
		}

		// re-init with the new branch so that we're in the same state as an existing branch. Probably a smarter way to achieve this.
		return r.init()
	}
}

func (repo *repo) shallowClone() error {
	_, err := repo.execGitCmd("clone", "--branch", repo.settings.Branch, "--single-branch", "--depth=1", repo.settings.URL, repo.workspaceDir)
	return err
}

func (repo *repo) shallowCloneDefaultBranch() error {
	_, err := repo.execGitCmd("clone", "--single-branch", "--depth=1", repo.settings.URL, repo.workspaceDir)
	return err
}

// createEmptyBranch creates an empty branch from the current branch
func (repo *repo) createEmptyBranch() error {
	_, err := repo.execGitCmd("checkout", "-b", repo.settings.Branch)
	if err != nil {
		return err
	}
	_, err = repo.execGitCmd("rm", "--ignore-unmatch", "-rf", ".")
	if err != nil {
		return err
	}
	_, err = repo.execGitCmd("commit", "--allow-empty", "-m", "Initial commit for Riser state branch")
	return err
}

// branchExists determines if the branch exists on the remote. Returns an error if the remote is invalid
func (repo *repo) branchExists() (exists bool, err error) {
	buffer, err := repo.execGitCmd("ls-remote", "--heads", repo.settings.URL, repo.settings.Branch)
	if err != nil {
		return false, err
	}

	exists = strings.Contains(buffer.String(), fmt.Sprintf("refs/heads/%s", repo.settings.Branch))
	return exists, nil
}

func (repo *repo) fetch() error {
	_, err := repo.execGitCmd("fetch", "-f", remoteName, repo.settings.Branch)
	return err
}

func (repo *repo) execGitCmd(args ...string) (stdOutAndStdErr *bytes.Buffer, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitExecTimeoutSeconds)
	defer cancel()
	cmd := repo.buildGitCmd(ctx, args...)
	stdOutAndStdErr, err = execWithContext(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("git %s", args))
	}

	return stdOutAndStdErr, nil
}

func (repo *repo) buildGitCmd(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repo.workspaceDir
	return cmd
}
