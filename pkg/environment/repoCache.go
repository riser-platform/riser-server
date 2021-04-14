package environment

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/riser-platform/riser-server/pkg/git"
)

type RepoCache struct {
	settings RepoSettings
	newFunc  func(git.RepoSettings) (git.Repo, error)
	cache    map[string]git.Repo
	sync     sync.Mutex
}

func NewBranchPerEnvRepoCache(settings RepoSettings) *RepoCache {
	return &RepoCache{
		settings: settings,
		newFunc:  git.InitRepoWorkspace,
		cache:    map[string]git.Repo{},
	}
}

func NewFakeRepoCache() *RepoCache {
	return &RepoCache{
		newFunc: func(git.RepoSettings) (git.Repo, error) {
			return &git.FakeRepo{}, nil
		},
		cache: map[string]git.Repo{},
	}
}

func (cache *RepoCache) GetRepo(envName string) (git.Repo, error) {
	return cache.getRepo(envName, newGitSettingsForEnv(envName, cache.settings))
}

func (cache *RepoCache) getRepo(envName string, settings git.RepoSettings) (git.Repo, error) {
	if envName == "" {
		return nil, errors.New("Environment name cannot be empty")
	}

	cache.sync.Lock()
	defer cache.sync.Unlock()

	if repo, ok := cache.cache[envName]; ok {
		return repo, nil
	}

	repo, err := cache.newFunc(settings)
	if err != nil {
		return nil, err
	}

	cache.cache[envName] = repo
	return repo, nil
}

func newGitSettingsForEnv(envName string, settings RepoSettings) git.RepoSettings {
	return git.RepoSettings{
		URL:         settings.URL,
		LocalGitDir: filepath.Join(settings.BaseGitDir, "/env/", envName),
		Branch:      envName,
	}
}
