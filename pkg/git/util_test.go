package git

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

func Test_isNoChangesErr_cleanTreeErr(t *testing.T) {
	result := isNoChangesErr(errors.New("Your branch is up to date with 'origin/main'.\n\nnothing to commit, working tree clean\n"))

	assert.True(t, result)
}

func Test_isNoChangesErr_false(t *testing.T) {
	result := isNoChangesErr(errors.New("foo"))

	assert.False(t, result)
}
