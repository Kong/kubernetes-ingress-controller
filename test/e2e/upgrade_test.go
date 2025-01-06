//go:build e2e_tests

package e2e

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	dockerimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

const (
	// upgradeTestFromTag is the tag of the previous version of the controller to upgrade from.
	upgradeTestFromTag = "v3.0.1"

	// dblessURLTemplate is the template of the URL to the all-in-one-dbless.yaml manifest with a placeholder for the tag.
	dblessURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/test/e2e/manifests/all-in-one-dbless.yaml"

	// 	postgresURLTemplate is the template of the URL to the all-in-one-postgres.yaml manifest with a placeholder for the tag.
	postgresURLTemplate = "https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/%s/test/e2e/manifests/all-in-one-postgres.yaml"
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
	oldManifestPath := dumpToTempFile(t, oldManifestResp.Body)

	skipIfNotTwoConsecutiveKongMinorVersions(t, oldManifestPath, testParams.toManifestPath)

	t.Log("configuring upgrade manifests test")
	ctx, env := setupE2ETest(t)

	t.Logf("deploying previous kong manifests: %s", testParams.fromManifestURL)
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

	t.Logf("deploying target version of kong manifests: %s", testParams.toManifestPath)
	deployments := ManifestDeploy{
		Path: testParams.toManifestPath,
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

// skipIfNotTwoConsecutiveKongMinorVersions skips the test if the old and new Kong versions are not two consecutive
// minor versions. This is necessary because Kong in DB-mode doesn't support skipping minor versions when upgrading.
// See the Gateway upgrade guide for details: https://docs.konghq.com/gateway/latest/upgrade/.
func skipIfNotTwoConsecutiveKongMinorVersions(
	t *testing.T,
	oldManifestPath string,
	newManifestPath string,
) {
	oldKongVersion := extractKongVersionFromManifest(t, oldManifestPath)

	var newKongVersion kong.Version
	if targetKongImage := testenv.KongImageTag(); targetKongImage != "" {
		// If the target Kong image is specified via environment variable, use it...
		newKongVersion = extractKongVersionFromDockerImage(t, targetKongImage)
	} else {
		// ...otherwise, use the version used in the new manifest.
		newKongVersion = extractKongVersionFromManifest(t, newManifestPath)
	}

	if oldKongVersion.Major() != newKongVersion.Major() ||
		oldKongVersion.Minor()+1 != newKongVersion.Minor() {
		t.Skipf("skipping upgrade test because the old and new Kong versions are not two consecutive minor versions: %s and %s",
			oldKongVersion, newKongVersion)
	}
}

var kongVersionRegex = regexp.MustCompile(`image: (kong:.*)`)

// extractKongVersionFromManifest extracts the Kong version from the manifest.
func extractKongVersionFromManifest(t *testing.T, manifestPath string) kong.Version {
	manifest, err := os.ReadFile(manifestPath)
	require.NoError(t, err)

	res := kongVersionRegex.FindStringSubmatch(string(manifest))
	require.NotEmpty(t, res)

	version := res[1]
	return extractKongVersionFromDockerImage(t, version)
}

// extractKongVersionFromDockerImage extracts the Kong version from the docker image by inspecting the image's env vars
// for the KONG_VERSION env var.
func extractKongVersionFromDockerImage(t *testing.T, image string) kong.Version {
	dockerc, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)

	ctx := context.Background()

	t.Log("negotiating docker API version")
	dockerc.NegotiateAPIVersion(ctx)

	t.Logf("pulling docker image %s to inspect it", image)
	_, err = dockerc.ImagePull(ctx, image, dockerimage.PullOptions{})
	require.NoError(t, err)

	t.Logf("inspecting docker image %s", image)
	var imageDetails types.ImageInspect
	// Retry because the image may not be available immediately after pulling it.
	require.Eventually(t, func() bool {
		var err error
		imageDetails, _, err = dockerc.ImageInspectWithRaw(ctx, image)
		if err != nil {
			t.Logf("failed to inspect docker image %s: %s", image, err)
			return false
		}
		return true
	}, time.Minute, time.Second)

	kongVersionEnv, ok := lo.Find(imageDetails.Config.Env, func(s string) bool {
		return strings.HasPrefix(s, "KONG_VERSION=")
	})
	require.True(t, ok, "KONG_VERSION env var not found in image %s", image)

	version, err := kong.ParseSemanticVersion(strings.TrimPrefix(kongVersionEnv, "KONG_VERSION="))
	require.NoError(t, err)
	t.Logf("parsed Kong version %s from docker image %s", version, image)

	return version
}
