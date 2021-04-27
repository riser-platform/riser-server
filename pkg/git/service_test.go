package git

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildGitCmd(t *testing.T) {
	r := repo{workspaceDir: "/tmp/git"}
	ctx := context.Background()
	cmd := r.buildGitCmd(ctx, "status")

	assert.Equal(t, []string{"git", "status"}, cmd.Args)
	assert.Equal(t, "/tmp/git", cmd.Dir)
}
