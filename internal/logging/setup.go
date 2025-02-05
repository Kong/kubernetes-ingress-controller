package logging

import (
	"fmt"
	"io"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	kongdbreconciler "github.com/kong/go-database-reconciler/pkg/cprint"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// SetupLoggers sets up the loggers for the controller manager.
func SetupLoggers(c managercfg.Config, output io.Writer) (logr.Logger, error) {
	zapBase, err := makeLogger(c.LogLevel, c.LogFormat, output)
	if err != nil {
		return logr.Logger{}, fmt.Errorf("failed to make logger: %w", err)
	}
	logger := zapr.NewLoggerWithOptions(zapBase, zapr.LogInfoLevel("v"))

	// It's specific for the Kong Ingress Controller.
	if c.LogLevel != "trace" && c.LogLevel != "debug" {
		// Disable deck's per-change diff output
		kongdbreconciler.DisableOutput = true
	}

	// Prevents controller-runtime from logging
	// [controller-runtime] log.SetLogger(...) was never called; logs will not be displayed.
	ctrllog.SetLogger(logger)

	return logger, nil
}
