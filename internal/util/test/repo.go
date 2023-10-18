package test

import (
	"path"
	"path/filepath"
	"runtime"
)

func getRepoRoot() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	d := filepath.Dir(path.Join(path.Dir(b), "../../")) // Number of ../ depends on the path of this file.
	return filepath.Abs(d)
}
