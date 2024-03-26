//go:build integration_tests

package isolated

import (
	"context"
	"testing"
	"time"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestKongLicense(t *testing.T) {
	f := features.
		New("essentials").
		WithLabel(testlabels.Kind, testlabels.KindKongLicense).
		Setup(SkipIfEnterpriseNotEnabled).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		Assess(
			"Expect No Licenses found before creating KongLicense resource",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				adminURL := GetAdminURLFromCtx(ctx)
				require.NotNil(t, adminURL, "Should get URL to access Kong gateway admin APIs from context")
				licenses, err := helpers.GetKongLicenses(ctx, adminURL, consts.KongTestPassword)
				require.NoError(t, err, "Expect No errors in listing all licenses in Kong gateway")
				require.Len(t, licenses, 0, "Expect No licenses in Kong gateway now")

				return ctx
			},
		).Assess(
		"Expect Licenses available when KongLicense created",
		func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			licenseString, err := ktfkong.GetLicenseJSONFromEnv()
			require.NoError(t, err)

			kongLicenseResource := &kongv1alpha1.KongLicense{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-license",
				},
				RawLicenseString: licenseString,
				Enabled:          true,
			}
			cluster := GetClusterFromCtx(ctx)
			kongClients, err := clientset.NewForConfig(cluster.Config())
			require.NoError(t, err, "Should get clientset to operate KongLicense with no errors")
			_, err = kongClients.ConfigurationV1alpha1().KongLicenses().Create(ctx, kongLicenseResource, metav1.CreateOptions{})
			require.NoError(t, err, "Should return no errors on creating KongLicense")

			require.Eventually(t, func() bool {
				kongLicenseResource, err = kongClients.ConfigurationV1alpha1().KongLicenses().Get(ctx, kongLicenseResource.Name, metav1.GetOptions{})
				require.NoError(t, err, "Should not return error on getting the latest status of KongLicense")
				if len(kongLicenseResource.Status.KongLicenseControllerStatuses) != 1 {
					return false
				}
				controllerStatus := kongLicenseResource.Status.KongLicenseControllerStatuses[0]
				return lo.ContainsBy(controllerStatus.Conditions, func(c metav1.Condition) bool {
					return c.Type == "Programmed" && c.Status == metav1.ConditionTrue
				})
			}, time.Minute, time.Second,
				"Expect Programmed condition in controller status to become 'True' in one minute")

			adminURL := GetAdminURLFromCtx(ctx)
			require.Eventually(t, func() bool {
				licenses, err := helpers.GetKongLicenses(ctx, adminURL, consts.KongTestPassword)
				require.NoError(t, err, "Expect No errors in listing all licenses in Kong gateway")
				return len(licenses) == 1
			},
				time.Minute, time.Second,
				"Expect 1 license found in Kong gateway")

			return ctx
		},
	).Teardown(featureTeardown())
	tenv.Test(t, f.Feature())
}
