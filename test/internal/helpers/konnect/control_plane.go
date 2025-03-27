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
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect/sdk"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

// CreateTestControlPlane creates a control plane to be used in tests. It returns the created control plane's ID.
// It also sets up a cleanup function for it to be deleted.
func CreateTestControlPlane(ctx context.Context, t *testing.T) string {
	t.Helper()

	sdk := sdk.New(accessToken(), serverURLOpt())

	var cpID string
	createRgErr := retry.Do(func() error {
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

		cpID = createResp.ControlPlane.ID
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second))
	require.NoError(t, createRgErr)

	t.Cleanup(func() {
		ctx = context.Background()
		t.Logf("deleting test Konnect Control Plane: %q", cpID)
		err := retry.Do(
			func() error { //nolint:contextcheck
				_, err := sdk.ControlPlanes.DeleteControlPlane(context.Background(), cpID)
				return err
			},
			retry.Attempts(5), retry.Delay(time.Second),
		)
		assert.NoErrorf(t, err, "failed to cleanup a control plane: %q", cpID)

		// Since Konnect authorization v2 supports cleanup of roles after control plane deleted, we do not need to delete them manually.
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

	c, err := adminapi.NewKongClientForKonnectControlPlane(managercfg.KonnectConfig{
		ControlPlaneID: cpID,
		Address:        konnectControlPlaneAdminAPIBaseURL(),
		TLSClient: managercfg.TLSClientConfig{
			Cert: cert,
			Key:  key,
		},
	})
	require.NoError(t, err)
	return c
}
