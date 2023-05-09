package license_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	konnectLicense "github.com/kong/kubernetes-ingress-controller/v2/internal/konnect/license"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/license"
)

type mockUpstreamClient struct {
	listResponse *konnectLicense.ListLicenseResponse
}

func (m *mockUpstreamClient) List(context.Context, int) (*konnectLicense.ListLicenseResponse, error) {
	return m.listResponse, nil
}

func TestAgent(t *testing.T) {
	expectedLicense := &konnectLicense.Item{
		License:   "test-license",
		UpdatedAt: 1234567890,
	}
	upstreamClient := &mockUpstreamClient{
		listResponse: &konnectLicense.ListLicenseResponse{
			Items: []*konnectLicense.Item{
				expectedLicense,
			},
		},
	}

	a := license.NewAgent(upstreamClient, logr.Discard())
	ctx := context.Background()
	go func() {
		err := a.Start(ctx)
		require.NoError(t, err)
	}()

	require.Eventually(t, func() bool {
		actualLicense := a.GetLicense()
		if actualLicense.Payload == nil {
			return false
		}
		return *actualLicense.Payload == expectedLicense.License
	}, time.Second*5, time.Millisecond)
}
