package util

import (
	"os"
	"path"
)

// EnsureDir ensures that a directory exists for a given file or directory path. Note that a directory path must have a trailing slash
func EnsureDir(dirOrFilePath string, fileMode os.FileMode) error {
	dir := path.Dir(dirOrFilePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, fileMode)
	}

	return err
}

func EnsureTrailingSlash(path string) string {
	if path[len(path)-1] != '/' {
		return path + "/"
	}

	return path
}
