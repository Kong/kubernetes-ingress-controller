package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
)

type CleanupFunc func(context.Context) error

// ExitOnErrWithCode is a helper function meant for us in the test.Main to simplify failing and exiting
// the tests under unrecoverable error conditions. It will also attempt to perform any cluster
// cleaning necessary before exiting.
func ExitOnErrWithCode(ctx context.Context, err error, exitCode int, fns ...CleanupFunc) {
	if err == nil {
		return
	}

	fmt.Printf("WARNING: failure occurred: %v\n", err)
	for _, fn := range fns {
		if clErr := fn(ctx); clErr != nil {
			err = errors.Join(err, fmt.Errorf("cleanup failed after test failure occurred CLEANUP_FAILURE=(%w)", clErr))
		}
	}

	fmt.Fprintf(os.Stderr, "Error: tests failed: %s\n", err)
	os.Exit(exitCode)
}

// ExitOnErr is a wrapper around exitOnErrorWithCode that defaults to using the ExitCodeEnvSetupFailed
// exit code. This function is meant for convenience to wrap errors in setup that are hard to predict.
func ExitOnErr(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ExitOnErrWithCode(ctx, err, consts.ExitCodeEnvSetupFailed)
}
