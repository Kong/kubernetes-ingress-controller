package license_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/license"
)

type mockUpstreamClient struct {
	listResponse *konnect.ListLicenseResponse
}

func (m *mockUpstreamClient) List(context.Context, int) (*konnect.ListLicenseResponse, error) {
	return m.listResponse, nil
}

func TestAgent(t *testing.T) {
	expectedLicense := &konnect.LicenseItem{
		License:   "test-license",
		UpdatedAt: 1234567890,
	}
	upstreamClient := &mockUpstreamClient{
		listResponse: &konnect.ListLicenseResponse{
			Items: []*konnect.LicenseItem{
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
	}, time.Second*5, time.Millisecond*100)
}
