package git

import (
	"errors"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	git "gopkg.in/src-d/go-git.v4"
)

func Test_handleGetRepository(t *testing.T) {
	repo := &git.Repository{}
	repoSettings := RepoSettings{
		URL: "https://tempuri.org",
	}
	getRepoFunc := func(s RepoSettings) (*git.Repository, error) {
		return repo, nil
	}

	result, err := handleGetRepository(repoSettings, getRepoFunc)

	require.Nil(t, err)
	assert.Equal(t, repo, result.repo)
	assert.Equal(t, repoSettings, result.repoSettings)
}

func Test_handleGetRepository_Error(t *testing.T) {
	expectedErr := errors.New("broke")
	getRepoFunc := func(s RepoSettings) (*git.Repository, error) {
		return nil, expectedErr
	}

	result, err := handleGetRepository(RepoSettings{}, getRepoFunc)

	require.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func Test_createCloneOptions(t *testing.T) {
	settings := RepoSettings{
		URL:    "https://tempuri.org/repo",
		Branch: "release",
	}

	result := createCloneOptions(settings)

	require.NotNil(t, result)
	assert.Equal(t, settings.URL, result.URL)
	assert.Equal(t, "refs/heads/release", result.ReferenceName.String())
	assert.Nil(t, result.Auth)
}

func Test_createCloneOptions_WithCreds(t *testing.T) {
	settings := RepoSettings{
		URL:      "https://tempuri.org/repo",
		Branch:   "release",
		Username: "myuser",
		Password: "mypass",
	}

	result := createCloneOptions(settings)

	require.NotNil(t, result)
	assert.Equal(t, settings.URL, result.URL)
	assert.Equal(t, "refs/heads/release", result.ReferenceName.String())
	assert.IsType(t, &http.BasicAuth{}, result.Auth)
	authResult := result.Auth.(*http.BasicAuth)
	assert.Equal(t, settings.Username, authResult.Username)
	assert.Equal(t, settings.Password, authResult.Password)
}
