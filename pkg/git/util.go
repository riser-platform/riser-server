package git

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
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
	// We require both username and password to be set in order to pass auth
	if settings.Username != "" && settings.Password != "" {
		parsedUrl.User = url.UserPassword(settings.Username, settings.Password)
	}
	return parsedUrl.String(), nil
}

func writeFiles(baseDir string, files []core.ResourceFile) error {
	for _, file := range files {
		fullFileName := filepath.Join(baseDir, file.Name)
		err := util.EnsureDir(fullFileName, workspaceFilePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fullFileName, file.Contents, workspaceFilePerm)
		if err != nil {
			return errors.Wrap(err, "error writing file")
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
