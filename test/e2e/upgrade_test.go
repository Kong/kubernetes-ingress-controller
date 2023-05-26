//go:build e2e_tests

package e2e

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const (
	// upgradeTestFromTag is the tag of the previous version of the controller to upgrade from.
	upgradeTestFromTag = "v2.9.3"

	// dblessURLTemplate is the template of the URL to the all-in-one-dbless.yaml manifest with a placeholder for the tag.
	dblessURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/deploy/single/all-in-one-dbless.yaml"

	// 	postgresURLTemplate is the template of the URL to the all-in-one-postgres.yaml manifest with a placeholder for the tag.
	postgresURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/deploy/single/all-in-one-postgres.yaml"
)

func TestDeployAndUpgradeAllInOneDBLESS(t *testing.T) {
	fromManifestURL := fmt.Sprintf(dblessURLTemplate, upgradeTestFromTag)
	toManifestPath := dblessPath
	testManifestsUpgrade(t, fromManifestURL, toManifestPath)
}

func TestDeployAndUpgradeAllInOnePostgres(t *testing.T) {
	fromManifestURL := fmt.Sprintf(postgresURLTemplate, upgradeTestFromTag)
	toManifestPath := postgresPath
	testManifestsUpgrade(t, fromManifestURL, toManifestPath)
}

func testManifestsUpgrade(t *testing.T, fromManifestURL string, toManifestPath string) {
	t.Parallel()

	httpClient := helpers.RetryableHTTPClient(helpers.DefaultHTTPClient())
	oldManifest, err := httpClient.Get(fromManifestURL)
	require.NoError(t, err)
	defer oldManifest.Body.Close()

	t.Log("configuring upgrade manifests test")
	ctx, env := setupE2ETest(t)

	t.Logf("deploying previous kong manifests: %s", fromManifestURL)
	deployKong(ctx, t, env, oldManifest.Body)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	t.Logf("deploying target version of kong manifests: %s", toManifestPath)
	manifest := getTestManifest(t, toManifestPath)
	deployKong(ctx, t, env, manifest)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}
