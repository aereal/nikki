package entrypoint

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/aereal/nikki/backend/log"
)

func provideDynamicServiceVersion(ctx context.Context) (log.ServiceVersion, error) {
	c := exec.CommandContext(ctx, "git", "describe", "--always", "--tags", "--dirty", "--abbrev=0")
	c.WaitDelay = time.Second * 1
	c.Cancel = func() error { return c.Process.Signal(os.Interrupt) }
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("command failed: stderr=%s: %w", stderr, err)
	}
	return log.ServiceVersion(bytes.TrimSpace(stdout.Bytes())), nil
}
