package konnect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	sdkkonnectgo "github.com/Kong/sdk-konnect-go"
	"github.com/Kong/sdk-konnect-go/retry"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"
	sdkkonnectops "github.com/Kong/sdk-konnect-go/models/operations"
	sdkkonnecterrs "github.com/Kong/sdk-konnect-go/models/sdkerrors"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/sdk"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// CreateTestControlPlane creates a control plane to be used in tests. It returns the created control plane's ID.
// It also sets up a cleanup function for it to be deleted.
func CreateTestControlPlane(ctx context.Context, t *testing.T) string {
	t.Helper()

	sdk := sdk.New(accessToken(), serverURLOpt(), sdkkonnectgo.WithRetryConfig(retry.Config{
		Strategy: "backoff",
		Backoff: &retry.BackoffStrategy{
			InitialInterval: 100,
			MaxInterval:     2000,
			Exponent:        1.2,
			MaxElapsedTime:  10000,
		},
	}))

	createResp, err := sdk.ControlPlanes.CreateControlPlane(ctx,
		sdkkonnectcomp.CreateControlPlaneRequest{
			Name:        uuid.NewString(),
			Description: lo.ToPtr(generateTestKonnectControlPlaneDescription(t)),
			Labels: map[string]string{
				test.KonnectControlPlaneLabelCreatedInTests: "true",
			},
			ClusterType: sdkkonnectcomp.CreateControlPlaneRequestClusterTypeClusterTypeK8SIngressController.ToPointer(),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotNil(t, createResp.ControlPlane)
	require.Equal(t, http.StatusCreated, createResp.GetStatusCode())
	require.NotNil(t, createResp.ControlPlane.ID)
	cpID := createResp.ControlPlane.ID

	t.Cleanup(func() {
		t.Logf("deleting test Konnect Control Plane: %q", cpID)
		_, err := sdk.ControlPlanes.DeleteControlPlane(ctx, cpID)
		assert.NoErrorf(t, err, "failed to cleanup a control plane: %q", cpID)

		me, err := sdk.Me.GetUsersMe(ctx,
			// NOTE: Otherwise we use prod server by default.
			// Related issue: https://github.com/Kong/sdk-konnect-go/issues/20
			sdkkonnectops.WithServerURL(test.KonnectServerURL()),
		)
		if !assert.NoError(t, err) {
			return
		}

		// We have to manually delete roles created for the control plane because Konnect doesn't do it automatically.
		// If we don't do it, we will eventually hit a problem with Konnect APIs answering our requests with 504s
		// because of a performance issue when there's too many roles for the account
		// (see https://konghq.atlassian.net/browse/TPS-1319).
		//
		// We can drop this once the automated cleanup is implemented on Konnect side:
		// https://konghq.atlassian.net/browse/TPS-1453.
		resp, err := sdk.Roles.ListUserRoles(ctx, *me.User.ID,
			&sdkkonnectops.ListUserRolesQueryParamFilter{},
			// Related issue: https://github.com/Kong/sdk-konnect-go/issues/20
			sdkkonnectops.WithServerURL(test.KonnectServerURL()),
		)
		require.NoErrorf(t, err, "failed to list control plane roles for cleanup: %q", cpID)

		for _, role := range resp.AssignedRoleCollection.Data {
			if role.EntityID == nil || role.ID == nil {
				continue
			}
			if *role.EntityID != cpID {
				continue
			}

			// Delete only roles created for the control plane.
			t.Logf("deleting test Konnect Control Plane role: %q", *role.ID)
			_, err := sdk.Roles.UsersRemoveRole(ctx, *me.User.ID, *role.ID,
				// Related issue: https://github.com/Kong/sdk-konnect-go/issues/20
				sdkkonnectops.WithServerURL(test.KonnectServerURL()),
			)
			notFoundErr := &sdkkonnecterrs.NotFoundError{}
			if !errors.As(err, &notFoundErr) {
				assert.NoErrorf(t, err, "failed to cleanup a control plane role: %q", *role.ID)
			}
		}
	})

	t.Logf("created test Konnect Control Plane: %q", cpID)
	return cpID
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
