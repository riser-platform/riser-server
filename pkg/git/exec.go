package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

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
