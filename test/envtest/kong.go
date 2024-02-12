package envtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

// runKongEnterpriseHandlingVaultValidationCorrectly runs a Kong container and ensures it didn't fall into a state that would
// cause it to return a 405 status code consistently when validating Vaults.
// It's a known bug on the Gateway side. Purpose of this function is to make sure that our tests are not affected by the bug.
// Once https://konghq.atlassian.net/browse/KAG-3699 is solved we should be able to drop the usage of this function.
func runKongEnterpriseHandlingVaultValidationCorrectly(ctx context.Context, t *testing.T) containers.Kong {
	// Get the Kong Gateway version to use for the test from `test_dependencies.yaml` file.
	gatewayTag, err := testenv.GetDependencyVersion("envtests.kong-ee")
	require.NoError(t, err)

	// Prepare the container config modifier to set the Kong Gateway version.
	withEnvtestsVersion := func(request *testcontainers.ContainerRequest) {
		request.Image = fmt.Sprintf("kong/kong-gateway:%s", gatewayTag)
	}

	var kongContainer containers.Kong
	require.Eventually(t, func() bool {
		kongContainer = containers.NewKong(ctx, t, withEnvtestsVersion)
		adminURL := kongContainer.AdminURL(ctx, t)

		kongClient, err := adminapi.NewKongClientForWorkspace(ctx, adminURL, "default", helpers.DefaultHTTPClient())
		if err != nil {
			t.Logf("Failed to create Kong client: %v", err)
			return false
		}

		_, _, err = kongClient.AdminAPIClient().Vaults.Validate(ctx, &kong.Vault{
			Name:        lo.ToPtr("env"),
			Description: lo.ToPtr("test-vault-description"),
			Prefix:      lo.ToPtr("test-vault-prefix"),
		})
		if err != nil {
			t.Logf("Vault validation endpoint malfunction discovered: %s. Retrying...", err)
			if err := kongContainer.Terminate(ctx); err != nil {
				t.Logf("Failed to terminate Kong container: %s", err)
			}
			return false
		}

		return true
	}, time.Second*60, time.Millisecond)

	return kongContainer
}
