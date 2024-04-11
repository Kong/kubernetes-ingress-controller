package helpers

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// GenerateKongBuilder returns a Kong KTF addon builder, a string slice
// of controller arguments needed to interact with the addon and an error.
func GenerateKongBuilder(_ context.Context) (*kong.Builder, []string, error) {
	kongbuilder := kong.NewBuilder().WithNamespace(consts.ControllerNamespace)
	extraControllerArgs := []string{}
	if testenv.KongEnterpriseEnabled() {
		licenseJSON, err := kong.GetLicenseJSONFromEnv()
		if err != nil {
			return nil, nil, err
		}
		kongbuilder = kongbuilder.WithProxyEnterpriseEnabled(licenseJSON)
		if testenv.DBMode() != testenv.DBModeOff {
			kongbuilder.WithProxyEnterpriseSuperAdminPassword(consts.KongTestPassword)
			extraControllerArgs = append(extraControllerArgs,
				fmt.Sprintf("--kong-admin-token=%s", consts.KongTestPassword),
				fmt.Sprintf("--kong-workspace=%s", consts.KongTestWorkspace),
			)
		}
	}

	if image, tag := testenv.KongImage(), testenv.KongTag(); image != "" && tag != "" {
		kongbuilder = kongbuilder.WithProxyImage(image, tag)
	} else if tag != "" || image != "" {
		return nil, nil, fmt.Errorf("when specifying TEST_KONG_IMAGE or TEST_KONG_TAG, both need to be provided")
	}

	if effectiveKongVersion := testenv.KongEffectiveVersion(); effectiveKongVersion != "" {
		kongbuilder = kongbuilder.WithAdditionalValue("image.effectiveSemver", effectiveKongVersion)
	}

	if user, pass := testenv.KongPullUsername(), testenv.KongPullPassword(); user != "" || pass != "" {
		if user == "" || pass == "" {
			return nil, nil, fmt.Errorf("TEST_KONG_PULL_USERNAME requires TEST_KONG_PULL_PASSWORD")
		}
		kongbuilder = kongbuilder.WithProxyImagePullSecret("", user, pass, "")
	}

	if testenv.DBMode() == testenv.DBModePostgres {
		kongbuilder = kongbuilder.WithPostgreSQL()
	}

	flavor := testenv.KongRouterFlavor()
	if len(flavor) == 0 {
		flavor = dpconf.RouterFlavorTraditional
	}
	kongbuilder = kongbuilder.WithProxyEnvVar("router_flavor", string(flavor))

	kongbuilder.WithControllerDisabled()
	kongbuilder.WithProxyAdminServiceTypeLoadBalancer()

	return kongbuilder, extraControllerArgs, nil
}
