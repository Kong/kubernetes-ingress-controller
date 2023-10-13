package test

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

func getRepoRoot(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	dir := strings.TrimSpace(filepath.Clean(string(out)))
	return dir, nil
}
