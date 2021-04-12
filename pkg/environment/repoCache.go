package environment

import (
	"path/filepath"
	"sync"

	"github.com/riser-platform/riser-server/pkg/git"
)

type RepoCache struct {
	settings RepoSettings
	newFunc  func(git.RepoSettings) (git.Repo, error)
	repos    map[string]git.Repo
	sync     sync.Mutex
}

func NewRepoCache(settings RepoSettings) *RepoCache {
	return &RepoCache{
		newFunc: git.NewRepo,
		repos:   map[string]git.Repo{},
	}
}

func (cache *RepoCache) GetRepo(envName string) (git.Repo, error) {
	return cache.getRepo(envName, newGitSettingsForEnv(envName, cache.settings))
}

func newFakeRepoCache() *RepoCache {
	return &RepoCache{
		newFunc: func(git.RepoSettings) (git.Repo, error) {
			return &git.FakeRepo{}, nil
		},
		repos: map[string]git.Repo{},
	}
}

func (cache *RepoCache) getRepo(envName string, settings git.RepoSettings) (git.Repo, error) {
	cache.sync.Lock()
	defer cache.sync.Unlock()

	if repo, ok := cache.repos[envName]; ok {
		return repo, nil
	}

	repo, err := cache.newFunc(settings)
	// TODO: test
	if err != nil {
		return nil, err
	}

	cache.repos[envName] = repo
	return repo, nil
}

func newGitSettingsForEnv(envName string, settings RepoSettings) git.RepoSettings {
	return git.RepoSettings{
		URL:         settings.URL,
		LocalGitDir: filepath.Join(settings.LocalGitDir, "/env/", envName),
		Branch:      envName,
	}
}
