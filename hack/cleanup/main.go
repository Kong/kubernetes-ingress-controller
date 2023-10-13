// This script cleans up orphaned GKE clusters and Konnect runtime
// groups that were created by the e2e tests (caued by e.g. unexpected
// crash that didn't allow a test's teardown to be completed correctly).
// It's meant to be installed as a cronjob and run repeatedly throughout
// the day to catch any orphaned resources: however tests should be trying to
// delete the resources they create themselves.
//
// A cluster is considered orphaned when all conditions are satisfied:
// 1. Its name begins with a predefined prefix (`gke-e2e-`).
// 2. It was created more than 1h ago.
//
// A control plane is considered orphaned when all conditions are satisfied:
// 1. It has a label `created_in_tests` with value `true`.
// 2. It was created more than 1h ago.
//
// Usage: `go run ./hack/cleanup [mode]`
// Where `mode` is one of:
// - `all` (default): clean up both GKE clusters and Konnect control planes
// - `gke`: clean up only GKE clusters
// - `konnect`: clean up only Konnect control planes
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"go.uber.org/zap"
)

const (
	konnectAccessTokenVar = "TEST_KONG_KONNECT_ACCESS_TOKEN" //nolint:gosec

	cleanupModeAll     = "all"
	cleanupModeGKE     = "gke"
	cleanupModeKonnect = "konnect"
)

var (
	gkeCreds           = os.Getenv(gke.GKECredsVar)
	gkeProject         = os.Getenv(gke.GKEProjectVar)
	gkeLocation        = os.Getenv(gke.GKELocationVar)
	konnectAccessToken = os.Getenv(konnectAccessTokenVar)
)

func main() {
	zaplog, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	log := zapr.NewLogger(zaplog)

	mode, err := getCleanupMode()
	if err != nil {
		log.Error(err, "error getting cleanup mode")
		os.Exit(1)
	}

	if err := validateVars(mode); err != nil {
		log.Error(err, "error validating vars")
		os.Exit(1)
	}

	cleanupFuncs := resolveCleanupFuncs(mode)
	ctx := context.Background()
	for _, f := range cleanupFuncs {
		if err := f(ctx, log); err != nil {
			log.Error(err, "error running cleanup function")
			os.Exit(1)
		}
	}
}

func getCleanupMode() (string, error) {
	if len(os.Args) < 2 {
		return cleanupModeAll, nil
	}

	switch os.Args[1] {
	case cleanupModeGKE:
	case cleanupModeKonnect:
	default:
		return "", fmt.Errorf("invalid cleanup mode: %s", os.Args[1])
	}

	return os.Args[1], nil
}

func resolveCleanupFuncs(mode string) []func(context.Context, logr.Logger) error {
	switch mode {
	case cleanupModeGKE:
		return []func(context.Context, logr.Logger) error{
			cleanupGKEClusters,
		}
	case cleanupModeKonnect:
		return []func(context.Context, logr.Logger) error{
			cleanupKonnectControlPlanes,
		}
	default:
		return []func(context.Context, logr.Logger) error{
			cleanupGKEClusters,
			cleanupKonnectControlPlanes,
		}
	}
}

func validateVars(mode string) error {
	switch mode {
	case cleanupModeGKE:
		return validateGKEVars()
	case cleanupModeKonnect:
		return validateKonnectVars()
	default:
		if err := validateGKEVars(); err != nil {
			return err
		}
		if err := validateKonnectVars(); err != nil {
			return err
		}
		return nil
	}
}

func validateKonnectVars() error {
	return notEmpty(konnectAccessTokenVar, konnectAccessToken)
}

func validateGKEVars() error {
	if err := notEmpty(gke.GKECredsVar, gkeCreds); err != nil {
		return err
	}
	if err := notEmpty(gke.GKEProjectVar, gkeProject); err != nil {
		return err
	}
	return notEmpty(gke.GKELocationVar, gkeLocation)
}

func notEmpty(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s was empty", name)
	}
	return nil
}
