//go:build conformance_tests

package conformance

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/gateway-api/conformance/apis/v1alpha1"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/yaml"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/metadata"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

func TestGatewayExperimentalConformance(t *testing.T) {
	if !shouldRunExperimentalConformance() {
		t.Skip("skipping experimental conformance tests")
	}
	// Experimental conformance tests passes only when expression routes are enabled
	// (KONG_TEST_EXPRESSION_ROUTES='true' set automatically in make test.conformance-experimental)
	// because it supports more features and it's the desired way to use Kong in the future.
	client, gatewayClassName := prepareEnvForGatewayConformanceTests(t)
	// Release is expected to have format v2.8.1 (when run on tagged commit) or v2.8.1-731-g4a32cf0b.
	_, err := semver.Parse(strings.TrimPrefix(metadata.Release, "v"))
	require.NoError(t, err)

	cSuite, err := suite.NewExperimentalConformanceTestSuite(
		suite.ExperimentalConformanceOptions{
			Options: suite.Options{
				Client:               client,
				GatewayClassName:     gatewayClassName,
				Debug:                true,
				CleanupBaseResources: !testenv.IsCI(),
				BaseManifests:        conformanceTestsBaseManifests,
				SupportedFeatures: sets.New(
					suite.SupportHTTPRouteQueryParamMatching,
					suite.SupportHTTPRouteMethodMatching,
					suite.SupportHTTPRouteResponseHeaderModification,
				),
			},
			ConformanceProfiles: sets.New(
				suite.HTTPConformanceProfileName,
			),
			Implementation: v1alpha1.Implementation{
				Organization: metadata.Organization,
				Project:      metadata.ProjectName,
				URL:          metadata.ProjectURL,
				Version:      metadata.Release,
				Contact: []string{
					path.Join(metadata.ProjectURL, "/issues/new/choose"),
				},
			},
		},
	)
	require.NoError(t, err)
	t.Log("starting the gateway conformance test suite")
	cSuite.Setup(t)
	// To work with individual tests only, you can disable the normal Run call and construct a slice containing a
	// single test only, e.g.:
	//
	//cSuite.Run(t, []suite.ConformanceTest{tests.HTTPRouteRedirectPortAndScheme})
	require.NoError(t, cSuite.Run(t, tests.ConformanceTests))

	const reportFileName = "kong-kubernetes-ingress-controller.yaml"
	t.Log("saving the gateway conformance test report to file:", reportFileName)
	report, err := cSuite.Report()
	require.NoError(t, err)
	rawReport, err := yaml.Marshal(report)
	require.NoError(t, err)
	// Save report in root of the repository, file name is in .gitignore.
	require.NoError(t, os.WriteFile("../../"+reportFileName, rawReport, 0o600))
}
