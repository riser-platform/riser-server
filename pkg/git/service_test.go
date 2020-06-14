package git

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildGitCmd(t *testing.T) {
	repoSettings := &RepoSettings{
		LocalGitDir: "/tmp/git",
	}
	r := repo{settings: repoSettings}
	ctx := context.Background()
	cmd := r.buildGitCmd(ctx, "status")

	assert.Equal(t, []string{"git", "status"}, cmd.Args)
	assert.Equal(t, "/tmp/git", cmd.Dir)
}
