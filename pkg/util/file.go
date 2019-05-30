package util

import (
	"os"
	"path"
)

func EnsureDir(dirOrFilePath string, fileMode os.FileMode) error {
	dir := path.Dir(dirOrFilePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, fileMode)
	}

	return err
}
