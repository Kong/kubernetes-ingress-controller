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
// A runtime group is considered orphaned when all conditions are satisfied:
// 1. It has a label `created_in_tests` with value `true`.
// 2. It was created more than 1h ago.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/gke"
	"github.com/sirupsen/logrus"
)

const (
	konnectAccessTokenVar = "TEST_KONG_KONNECT_ACCESS_TOKEN" //nolint:gosec
)

var (
	gkeCreds           = os.Getenv(gke.GKECredsVar)
	gkeProject         = os.Getenv(gke.GKEProjectVar)
	gkeLocation        = os.Getenv(gke.GKELocationVar)
	konnectAccessToken = os.Getenv(konnectAccessTokenVar)
	log                = logrus.New()
)

func main() {
	validateVars()

	ctx := context.Background()
	if err := cleanupGKEClusters(ctx); err != nil {
		log.Errorf("error cleaning up GKE clusters: %v\n", err)
		os.Exit(1)
	}

	if err := cleanupKonnectRuntimeGroups(ctx); err != nil {
		log.Errorf("error cleaning up Konnect runtime groups: %v\n", err)
		os.Exit(1)
	}
}

func validateVars() {
	mustNotBeEmpty(gke.GKECredsVar, gkeCreds)
	mustNotBeEmpty(gke.GKEProjectVar, gkeProject)
	mustNotBeEmpty(gke.GKELocationVar, gkeLocation)
	mustNotBeEmpty(konnectAccessTokenVar, konnectAccessToken)
}

func mustNotBeEmpty(name, value string) {
	if value == "" {
		panic(fmt.Sprintf("%s was empty", name))
	}
}
