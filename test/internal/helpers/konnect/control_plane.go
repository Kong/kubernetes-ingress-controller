package konnect

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/roles"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/sdk"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// CreateTestControlPlane creates a control plane to be used in tests. It returns the created control plane's ID.
// It also sets up a cleanup function for it to be deleted.
func CreateTestControlPlane(ctx context.Context, t *testing.T) string {
	t.Helper()

	sdk := sdk.SDK(accessToken(), serverURLOpt())

	rolesClient := roles.NewClient(
		helpers.RetryableHTTPClient(helpers.DefaultHTTPClient()),
		konnectRolesBaseURL(),
		accessToken(),
	)

	var rgID string
	createRgErr := retry.Do(func() error {
		createResp, err := sdk.ControlPlanes.CreateControlPlane(ctx,
			sdkkonnectcomp.CreateControlPlaneRequest{
				Name:        uuid.NewString(),
				Description: lo.ToPtr(generateTestKonnectControlPlaneDescription(t)),
				Labels: map[string]string{
					"created_in_tests": "true",
				},
				ClusterType: sdkkonnectcomp.CreateControlPlaneRequestClusterTypeClusterTypeK8SIngressController.ToPointer(),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create control plane: %w", err)
		}
		if createResp == nil || createResp.ControlPlane == nil {
			return fmt.Errorf("failed to create control plane: response is nil, status code %d, response: %v",
				createResp.GetStatusCode(), createResp)
		}

		if createResp.GetStatusCode() != http.StatusCreated {
			body, err := io.ReadAll(createResp.RawResponse.Body)
			if err != nil {
				body = []byte(err.Error())
			}
			return fmt.Errorf("failed to create RG: code %d, message %s", createResp.GetStatusCode(), body)
		}
		if createResp.ControlPlane == nil || createResp.ControlPlane.ID == "" {
			return errors.New("No control plane ID in response")
		}

		rgID = createResp.ControlPlane.ID
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second))
	require.NoError(t, createRgErr)

	t.Cleanup(func() {
		t.Logf("deleting test Konnect Control Plane: %q", rgID)
		err := retry.Do(
			func() error {
				_, err := sdk.ControlPlanes.DeleteControlPlane(ctx, rgID)
				return err
			},
			retry.Attempts(5), retry.Delay(time.Second),
		)
		assert.NoErrorf(t, err, "failed to cleanup a control plane: %q", rgID)

		// We have to manually delete roles created for the control plane because Konnect doesn't do it automatically.
		// If we don't do it, we will eventually hit a problem with Konnect APIs answering our requests with 504s
		// because of a performance issue when there's too many roles for the account
		// (see https://konghq.atlassian.net/browse/TPS-1319).
		//
		// We can drop this once the automated cleanup is implemented on Konnect side:
		// https://konghq.atlassian.net/browse/TPS-1453.
		rgRoles, err := rolesClient.ListControlPlanesRoles(ctx)
		require.NoErrorf(t, err, "failed to list control plane roles for cleanup: %q", rgID)
		for _, role := range rgRoles {
			if role.EntityID == rgID { // Delete only roles created for the control plane.
				t.Logf("deleting test Konnect Control Plane role: %q", role.ID)
				err := rolesClient.DeleteRole(ctx, role.ID)
				assert.NoErrorf(t, err, "failed to cleanup a control plane role: %q", role.ID)
			}
		}
	})

	t.Logf("created test Konnect Control Plane: %q", rgID)
	return rgID
}

func generateTestKonnectControlPlaneDescription(t *testing.T) string {
	t.Helper()

	desc := fmt.Sprintf("control plane for test %s", t.Name())
	if testenv.GithubServerURL() != "" && testenv.GithubRepo() != "" && testenv.GithubRunID() != "" {
		githubRunURL := fmt.Sprintf("%s/%s/actions/runs/%s",
			testenv.GithubServerURL(), testenv.GithubRepo(), testenv.GithubRunID())
		desc += ", github workflow run " + githubRunURL
	}

	return desc
}

// CreateKonnectAdminAPIClient creates an *kong.Client that will communicate with Konnect Control Plane's Admin API.
func CreateKonnectAdminAPIClient(t *testing.T, cpID, cert, key string) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongClientForKonnectControlPlane(adminapi.KonnectConfig{
		ControlPlaneID: cpID,
		Address:        konnectControlPlaneAdminAPIBaseURL(),
		TLSClient: adminapi.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	})
	require.NoError(t, err)
	return c
}
