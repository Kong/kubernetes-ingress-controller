package konnect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	cp "github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/controlplanes"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/roles"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// CreateTestControlPlane creates a control plane to be used in tests. It returns the created control plane's ID.
// It also sets up a cleanup function for it to be deleted.
func CreateTestControlPlane(ctx context.Context, t *testing.T) string {
	t.Helper()
	rgClient, err := cp.NewClientWithResponses(konnectControlPlanesBaseURL, cp.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+accessToken())
			return nil
		}),
	)
	require.NoError(t, err)
	rolesClient := roles.NewClient(
		helpers.RetryableHTTPClient(helpers.DefaultHTTPClient()),
		konnectRolesBaseURL,
		accessToken(),
	)

	var rgID uuid.UUID
	createRgErr := retry.Do(func() error {
		rgName := uuid.NewString()
		createRgResp, err := rgClient.CreateControlPlaneWithResponse(ctx, cp.CreateControlPlaneRequest{
			Description: lo.ToPtr(generateTestKonnectControlPlaneDescription(t)),
			Labels: &cp.Labels{
				"created_in_tests": "true",
			},
			Name:        rgName,
			ClusterType: cp.ClusterTypeKubernetesIngressController,
		})
		if err != nil {
			return fmt.Errorf("failed to create control plane: %w", err)
		}
		if createRgResp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("failed to create RG: code %d, message %s", createRgResp.StatusCode(), string(createRgResp.Body))
		}
		if createRgResp.JSON201 == nil || createRgResp.JSON201.Id == nil {
			return errors.New("No control plane ID in response")
		}

		rgID = *createRgResp.JSON201.Id
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second))
	require.NoError(t, createRgErr)

	t.Cleanup(func() {
		t.Logf("deleting test Konnect Control Plane: %q", rgID)
		err := retry.Do(
			func() error {
				_, err := rgClient.DeleteControlPlaneWithResponse(ctx, rgID)
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
			if role.EntityID == rgID.String() { // Delete only roles created for the control plane.
				t.Logf("deleting test Konnect Control Plane role: %q", role.ID)
				err := rolesClient.DeleteRole(ctx, role.ID)
				assert.NoErrorf(t, err, "failed to cleanup a control plane role: %q", role.ID)
			}
		}
	})

	t.Logf("created test Konnect Control Plane: %q", rgID.String())
	return rgID.String()
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
		Address:        konnectControlPlaneAdminAPIBaseURL,
		TLSClient: adminapi.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	})
	require.NoError(t, err)
	return c
}
