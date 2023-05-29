//go:build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	testManifestsUpgrade(t, fromManifestURL, toManifestPath, nil)
}

func TestDeployAndUpgradeAllInOnePostgres(t *testing.T) {
	fromManifestURL := fmt.Sprintf(postgresURLTemplate, upgradeTestFromTag)
	toManifestPath := postgresPath

	// Injecting a beforeUpgradeHook to delete the old migrations job before the upgrade. This is necessary because it's
	// not allowed to modify the existing job's spec.
	beforeUpgradeHook := func(ctx context.Context, t *testing.T, env environments.Environment) {
		err := env.Cluster().Client().BatchV1().Jobs(namespace).Delete(ctx, migrationsJobName, metav1.DeleteOptions{
			PropagationPolicy: lo.ToPtr(metav1.DeletePropagationBackground),
		})
		require.NoError(t, err, "failed to delete old migrations job before upgrade")
	}
	testManifestsUpgrade(t, fromManifestURL, toManifestPath, beforeUpgradeHook)
}

// beforeUpgradeFn is a function that is run before the upgrade to clean up any resources that may interfere with the upgrade.
type beforeUpgradeFn func(ctx context.Context, t *testing.T, env environments.Environment)

func testManifestsUpgrade(
	t *testing.T,
	fromManifestURL string,
	toManifestPath string,
	beforeUpgradeHook beforeUpgradeFn,
) {
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

	if beforeUpgradeHook != nil {
		t.Log("running before upgrade hook")
		beforeUpgradeHook(ctx, t, env)
	}

	t.Logf("deploying target version of kong manifests: %s", toManifestPath)
	manifest := getTestManifest(t, toManifestPath)
	deployKong(ctx, t, env, manifest)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
}
