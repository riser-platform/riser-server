package git

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_execWithContext(t *testing.T) {
	ctx := context.Background()
	cmd := exec.Command("echo", "arg0", "arg1")
	result, err := execWithContext(ctx, cmd)

	assert.NoError(t, err)
	assert.Equal(t, "arg0 arg1\n", result.String())
}

func Test_execWithContext_exitCodeErr(t *testing.T) {
	ctx := context.Background()
	cmd := exec.Command("sh", "-c", "echo test;exit 113")
	result, err := execWithContext(ctx, cmd)

	assert.Equal(t, "error: exit status 113, output:\ntest\n", err.Error())
	assert.Equal(t, "test\n", result.String())
}

func Test_execWithContext_exitCodeErrNoBytes(t *testing.T) {
	ctx := context.Background()
	cmd := exec.Command("sh", "-c", "exit 113")
	result, err := execWithContext(ctx, cmd)

	assert.Len(t, result.Bytes(), 0)
	assert.Equal(t, "error: exit status 113", err.Error())
}

func Test_execWithContext_timeout(t *testing.T) {
	timeout := 1 * time.Microsecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("sleep", "1")
	_, err := execWithContext(ctx, cmd)

	assert.Equal(t, "command aborted: /bin/sleep [sleep 1]: context deadline exceeded", err.Error())
}

func Test_execWithContext_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd := exec.Command("sleep", "1")
	_, err := execWithContext(ctx, cmd)

	assert.Equal(t, "command aborted: /bin/sleep [sleep 1]: context canceled", err.Error())
}
