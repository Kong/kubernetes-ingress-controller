//go:build e2e_tests

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const (
	// upgradeTestFromTag is the tag of the previous version of the controller to upgrade from.
	upgradeTestFromTag = "v2.9.4"

	// dblessURLTemplate is the template of the URL to the all-in-one-dbless.yaml manifest with a placeholder for the tag.
	dblessURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/deploy/single/all-in-one-dbless.yaml"

	// 	postgresURLTemplate is the template of the URL to the all-in-one-postgres.yaml manifest with a placeholder for the tag.
	postgresURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/deploy/single/all-in-one-postgres.yaml"
)

func TestDeployAndUpgradeAllInOneDBLESS(t *testing.T) {
	testManifestsUpgrade(t, manifestsUpgradeTestParams{
		fromManifestURL: fmt.Sprintf(dblessURLTemplate, upgradeTestFromTag),
		toManifestPath:  dblessPath,
	})
}

func TestDeployAndUpgradeAllInOnePostgres(t *testing.T) {
	testManifestsUpgrade(t, manifestsUpgradeTestParams{
		fromManifestURL:   fmt.Sprintf(postgresURLTemplate, upgradeTestFromTag),
		toManifestPath:    postgresPath,
		beforeUpgradeHook: postgresBeforeUpgradeHook,
	})
}

func postgresBeforeUpgradeHook(ctx context.Context, t *testing.T, env environments.Environment) {
	// Injecting a beforeUpgradeHook to delete the old migrations job before the upgrade. This is necessary because it's
	// not allowed to modify the existing job's spec.
	err := env.Cluster().Client().BatchV1().Jobs(namespace).Delete(ctx, migrationsJobName, metav1.DeleteOptions{
		PropagationPolicy: lo.ToPtr(metav1.DeletePropagationBackground),
	})
	require.NoError(t, err, "failed to delete old migrations job before upgrade")
}

type beforeUpgradeFn func(ctx context.Context, t *testing.T, env environments.Environment)

type manifestsUpgradeTestParams struct {
	// fromManifestURL is the URL to the manifest to deploy before the upgrade.
	fromManifestURL string

	// toManifestPath is the path to the manifest to deploy after the upgrade.
	toManifestPath string

	// beforeUpgradeHook is a function that is run before the upgrade to clean up any resources that may interfere with the upgrade.
	beforeUpgradeHook beforeUpgradeFn

	// controllerFeatureGates contains feature gates to enable on the controller during the upgrade (e.g. "FillID=true").
	controllerFeatureGates string
}

func testManifestsUpgrade(
	t *testing.T,
	testParams manifestsUpgradeTestParams,
) {
	httpClient := helpers.RetryableHTTPClient(helpers.DefaultHTTPClient())
	oldManifestResp, err := httpClient.Get(testParams.fromManifestURL)
	require.NoError(t, err)
	defer oldManifestResp.Body.Close()

	t.Log("configuring upgrade manifests test")
	ctx, env := setupE2ETest(t)

	t.Logf("deploying previous kong manifests: %s", testParams.fromManifestURL)
	oldManifestPath := dumpToTempFile(t, oldManifestResp.Body)
	ManifestDeploy{
		Path:            oldManifestPath,
		SkipTestPatches: true,
	}.Run(ctx, t, env)

	t.Log("running ingress tests to verify all-in-one deployed ingress controller and proxy are functional")
	ingress := deployIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)
	verifyIngressWithEchoBackends(ctx, t, env, numberOfEchoBackends)

	if hook := testParams.beforeUpgradeHook; hook != nil {
		t.Log("running before upgrade hook")
		hook(ctx, t, env)
	}

	// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/4973
	// This router flavor substitution is a workaround for the fact that KIC doesn't
	// support configuring Gateways of different router flavors at the same time.
	// Since pre 3.0 the default router flavor was "traditional", use it here during
	// the upgrade.
	// When 4973 is fixed, this will go away as those new tests will use a different
	// testing strategy.
	// Using expressions router in earlier version (e.g. 2.12) was causing
	// configuration translation errors.
	newManifestB, err := os.ReadFile(testParams.toManifestPath)
	require.NoError(t, err)
	newManifestB = bytes.ReplaceAll(newManifestB, []byte("value: expressions"), []byte("value: traditional"))
	newManifestPath := dumpToTempFile(t, bytes.NewReader(newManifestB))

	t.Logf("deploying target version of kong manifests: %s", testParams.toManifestPath)
	deployments := ManifestDeploy{
		Path: newManifestPath,
		// Do not skip test patches - we want to verify that upgrade works with an image override in target manifest.
		SkipTestPatches: false,
	}.Run(ctx, t, env)

	if featureGates := testParams.controllerFeatureGates; featureGates != "" {
		t.Logf("setting environment variables for controller feature gates: %s", featureGates)
		kubeconfig := getTemporaryKubeconfig(t, env)
		require.NoError(t, setEnv(setEnvParams{
			kubeCfgPath:   kubeconfig,
			namespace:     namespace,
			target:        fmt.Sprintf("deployment/%s", controllerDeploymentName),
			containerName: controllerContainerName,
			variableName:  "CONTROLLER_FEATURE_GATES",
			value:         featureGates,
		}))
		helpers.WaitForDeploymentRollout(ctx, t, env.Cluster(), deployments.ControllerNN.Namespace, deployments.ControllerNN.Name)
	}

	t.Log("creating new ingress with new path /echo-new")
	newIngress := ingress.DeepCopy()
	newIngress.Name = "echo-new"
	const newPath = "/echo-new"
	newIngress.Spec.Rules[0].HTTP.Paths[0].Path = newPath
	_, err = env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Create(ctx, newIngress, metav1.CreateOptions{})
	require.NoError(t, err)

	verifyIngressWithEchoBackendsPath(ctx, t, env, numberOfEchoBackends, newPath)
}
