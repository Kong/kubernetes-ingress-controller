package helpers

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"

	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

// GenerateKongBuilder returns a Kong KTF addon builder, a string slice
// of controller arguments needed to interact with the addon and an error.
func GenerateKongBuilder(ctx context.Context) (*kong.Builder, []string, error) {
	kongbuilder := kong.NewBuilder()
	extraControllerArgs := []string{}
	if testenv.KongEnterpriseEnabled() {
		licenseJSON, err := kong.GetLicenseJSONFromEnv()
		if err != nil {
			return nil, nil, err
		}
		extraControllerArgs = append(extraControllerArgs,
			fmt.Sprintf("--kong-admin-token=%s", consts.KongTestPassword),
			"--kong-workspace=notdefault",
		)
		kongbuilder = kongbuilder.WithProxyEnterpriseEnabled(licenseJSON).
			WithProxyEnterpriseSuperAdminPassword(consts.KongTestPassword).
			WithProxyAdminServiceTypeLoadBalancer()
	}

	if image, tag := testenv.KongImage(), testenv.KongTag(); image != "" {
		if tag == "" {
			return nil, nil, fmt.Errorf("TEST_KONG_IMAGE requires TEST_KONG_TAG")
		}
		kongbuilder = kongbuilder.WithProxyImage(image, tag)
	}

	if user, pass := testenv.KongPullUsername(), testenv.KongPullPassword(); user != "" || pass != "" {
		if user == "" || pass == "" {
			return nil, nil, fmt.Errorf("TEST_KONG_PULL_USERNAME requires TEST_KONG_PULL_PASSWORD")
		}
		kongbuilder = kongbuilder.WithProxyImagePullSecret("", user, pass, "")
	}

	if testenv.DBMode() == "postgres" {
		kongbuilder = kongbuilder.WithPostgreSQL()
	}

	flavor := testenv.KongRouterFlavor()
	if flavor == "" {
		flavor = "traditional"
	}
	kongbuilder = kongbuilder.WithProxyEnvVar("router_flavor", flavor)

	kongbuilder.WithControllerDisabled()

	return kongbuilder, extraControllerArgs, nil
}
