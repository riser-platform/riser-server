package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/pkg/core"
)

const (
	workspaceFilePerm = 0755
)

func processFiles(baseDir string, files []core.ResourceFile) error {
	for _, file := range files {
		fullFileName := filepath.Join(baseDir, file.Name)
		if file.Delete {
			err := os.RemoveAll(fullFileName)
			if err != nil {
				return errors.Wrap(err, "error deleting file or directory")
			}
		} else {
			err := util.EnsureDir(fullFileName, workspaceFilePerm)
			if err != nil {
				return err
			}

			err = os.WriteFile(fullFileName, file.Contents, workspaceFilePerm)
			if err != nil {
				return errors.Wrap(err, "error writing file")
			}
		}
	}

	return nil
}

func isNoChangesErr(err error) bool {
	return strings.Contains(err.Error(), "working tree clean")
}
