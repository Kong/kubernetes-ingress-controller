package helpers

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	kongsemver "github.com/kong/semver/v4"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func getKongVersion() (semver.Version, error) {
	kongVersion, err := semver.Parse(testenv.KongEffectiveVersion())
	if err == nil {
		return kongVersion, nil
	}
	// Use kong/semver to parse the version string in the TEST_KONG_TAG in case it is a four-digit version like "3.15.0.0".
	kongsemverVersion, err := kongsemver.Parse(testenv.KongTag())
	if err != nil {
		return semver.Version{}, fmt.Errorf("could not parse Kong version from TEST_KONG_EFFECTIVE_VERSION or TEST_KONG_TAG: %w", err)
	}
	return semver.Version{
		Major: kongsemverVersion.Major,
		Minor: kongsemverVersion.Minor,
		Patch: kongsemverVersion.Patch,
	}, nil
}

// GenerateKongBuilder returns a Kong KTF addon builder, a string slice
// of controller arguments needed to interact with the addon and an error.
func GenerateKongBuilder(_ context.Context) (*kong.Builder, []string, error) {
	kongVersion, err := getKongVersion()
	if err != nil {
		return nil, nil, fmt.Errorf("could not determine Kong version: %w", err)
	}

	kongbuilder := kong.NewBuilder().WithNamespace(consts.ControllerNamespace)
	extraControllerArgs := []string{}
	if testenv.KongEnterpriseEnabled() || kongVersion.GTE(consts.ForceLicenseVersionCutoff) {
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
		// Use bitnamilegacy/postgresql since the bitnami/postgresql repository is gone.
		kongbuilder = kongbuilder.WithAdditionalValue("postgresql.image.repository", "bitnamilegacy/postgresql")
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

// GenerateKongBuilderWithController generates a Kong builder that installs both KIC and Kong gateway.
func GenerateKongBuilderWithController() (*kong.Builder, error) {
	kongbuilder := kong.NewBuilder().WithNamespace(consts.ControllerNamespace)

	if testenv.KongEnterpriseEnabled() {
		licenseJSON, err := kong.GetLicenseJSONFromEnv()
		if err != nil {
			return nil, err
		}
		kongbuilder = kongbuilder.WithProxyEnterpriseEnabled(licenseJSON)
		if testenv.DBMode() != testenv.DBModeOff {
			kongbuilder.WithProxyEnterpriseSuperAdminPassword(consts.KongTestPassword)
		}
	}

	if image, tag := testenv.KongImage(), testenv.KongTag(); image != "" && tag != "" {
		kongbuilder = kongbuilder.WithProxyImage(image, tag)
	} else if tag != "" || image != "" {
		return nil, fmt.Errorf("when specifying TEST_KONG_IMAGE or TEST_KONG_TAG, both need to be provided")
	}

	if effectiveKongVersion := testenv.KongEffectiveVersion(); effectiveKongVersion != "" {
		kongbuilder = kongbuilder.WithAdditionalValue("image.effectiveSemver", effectiveKongVersion)
	}

	if user, pass := testenv.KongPullUsername(), testenv.KongPullPassword(); user != "" || pass != "" {
		if user == "" || pass == "" {
			return nil, fmt.Errorf("TEST_KONG_PULL_USERNAME requires TEST_KONG_PULL_PASSWORD")
		}
		kongbuilder = kongbuilder.WithProxyImagePullSecret("", user, pass, "")
	}

	if testenv.DBMode() == testenv.DBModePostgres {
		kongbuilder = kongbuilder.WithPostgreSQL()
		// Use bitnamilegacy/postgresql since the bitnami/postgresql repository is gone.
		kongbuilder = kongbuilder.WithAdditionalValue("postgresql.image.repository", "bitnamilegacy/postgresql")
	}

	flavor := testenv.KongRouterFlavor()
	if len(flavor) == 0 {
		flavor = dpconf.RouterFlavorTraditional
	}
	kongbuilder = kongbuilder.WithProxyEnvVar("router_flavor", string(flavor))

	kongbuilder.WithProxyAdminServiceTypeLoadBalancer()

	return kongbuilder, nil
}
