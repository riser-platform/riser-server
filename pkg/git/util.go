package git

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/pkg/core"
)

const (
	workspaceFilePerm = 0755
)

// formatUrlWithAuth returns URL with auth info if specified
func formatUrlWithAuth(settings RepoSettings) (string, error) {
	parsedUrl, err := url.Parse(settings.URL)
	if err != nil {
		return "", err
	}

	if settings.Username != "" && settings.Password != "" {
		parsedUrl.User = url.UserPassword(settings.Username, settings.Password)
	} else if settings.Username != "" {
		// Support username only (e.g. https://<token>@domain) style auth
		parsedUrl.User = url.User(settings.Username)
	}

	return parsedUrl.String(), nil
}

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

			err = ioutil.WriteFile(fullFileName, file.Contents, workspaceFilePerm)
			if err != nil {
				return errors.Wrap(err, "error writing file")
			}
		}
	}

	return nil
}

func execWithContext(ctx context.Context, cmd *exec.Cmd) (stdOutAndStdErr *bytes.Buffer, err error) {
	stdOutAndStdErr = &bytes.Buffer{}
	cmd.Stdout = stdOutAndStdErr
	cmd.Stderr = stdOutAndStdErr
	err = cmd.Run()

	if err != nil {
		if len(stdOutAndStdErr.Bytes()) > 0 {
			err = fmt.Errorf("error: %v, output:\n%s", err, stdOutAndStdErr.String())
		} else {
			err = fmt.Errorf("error: %v", err)
		}
	}

	if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
		err = errors.Wrap(ctx.Err(), fmt.Sprintf("command aborted: %s %v", cmd.Path, cmd.Args))
	}

	return stdOutAndStdErr, err
}

func isNoChangesErr(err error) bool {
	return strings.Contains(err.Error(), "working tree clean")
}
