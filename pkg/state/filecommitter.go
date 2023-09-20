package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

type FileCommitter struct {
	basePath string
}

func NewFileCommitter(basePath string) *FileCommitter {
	return &FileCommitter{basePath}
}

func (committer *FileCommitter) Commit(message string, files []core.ResourceFile) error {
	for _, file := range files {
		fullpath := filepath.Join(committer.basePath, file.Name)
		err := util.EnsureDir(fullpath, 0755)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error creating directory for file %q", fullpath))
		}
		err = os.WriteFile(fullpath, file.Contents, 0644)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error writing file %q", fullpath))
		}
	}
	return nil
}
