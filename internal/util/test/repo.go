package test

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

func getRepoRoot() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get repo root: runtime.Caller(0) failed")
	}
	d := filepath.Dir(path.Join(path.Dir(b), "../../")) // Number of ../ depends on the path of this file.
	return filepath.Abs(d)
}
