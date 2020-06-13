package git

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

func Test_processFiles(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "riser-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := core.ResourceFile{
		Name:     "nested/test01",
		Contents: []byte("contents"),
	}

	err = processFiles(dir, []core.ResourceFile{file})

	assert.NoError(t, err)
	contents, err := ioutil.ReadFile(filepath.Join(dir, file.Name))
	assert.NoError(t, err)
	assert.EqualValues(t, "contents", contents)
}

func Test_processFiles_deletes(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "riser-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := core.ResourceFile{
		Name:     "nested/test01",
		Contents: []byte("contents"),
	}

	err = processFiles(dir, []core.ResourceFile{file})
	assert.NoError(t, err)

	// Make sure we can delete the whole folder
	file.Name = "nested"
	file.Delete = true

	err = processFiles(dir, []core.ResourceFile{file})
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "nested/"))
	assert.True(t, os.IsNotExist(err))
}

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

func Test_isNoChangesErr_cleanTreeErr(t *testing.T) {
	result := isNoChangesErr(errors.New("Your branch is up to date with 'origin/main'.\n\nnothing to commit, working tree clean\n"))

	assert.True(t, result)
}

func Test_isNoChangesErr_false(t *testing.T) {
	result := isNoChangesErr(errors.New("foo"))

	assert.False(t, result)
}
