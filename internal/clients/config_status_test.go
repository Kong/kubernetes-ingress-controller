package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
)

func TestChannelConfigNotifier(t *testing.T) {
	n := clients.NewChannelConfigNotifier(logr.Discard())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := n.SubscribeGatewayConfigStatus()

	// Call NotifyConfigStatus 5 times to make sure that the method is non-blocking.
	for i := 0; i < 5; i++ {
		n.NotifyGatewayConfigStatus(ctx, clients.GatewayConfigApplyStatus{})
	}

	for i := 0; i < 5; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for config status i=%d", i)
		}
	}
}

func TestCalculateConfigStatus(t *testing.T) {
	testCases := []struct {
		name string

		gatewayFailure      bool
		konnectFailure      bool
		translationFailures bool

		expectedConfigStatus clients.ConfigStatus
	}{
		{
			name:                 "success",
			expectedConfigStatus: clients.ConfigStatusOK,
		},
		{
			name:                 "gateway failure",
			gatewayFailure:       true,
			expectedConfigStatus: clients.ConfigStatusApplyFailed,
		},
		{
			name:                 "translation failures",
			translationFailures:  true,
			expectedConfigStatus: clients.ConfigStatusTranslationErrorHappened,
		},
		{
			name:                 "konnect failure",
			konnectFailure:       true,
			expectedConfigStatus: clients.ConfigStatusOKKonnectApplyFailed,
		},
		{
			name:                 "both gateway and konnect failure",
			gatewayFailure:       true,
			konnectFailure:       true,
			expectedConfigStatus: clients.ConfigStatusApplyFailedKonnectApplyFailed,
		},
		{
			name:                 "translation failures and konnect failure",
			translationFailures:  true,
			konnectFailure:       true,
			expectedConfigStatus: clients.ConfigStatusTranslationErrorHappenedKonnectApplyFailed,
		},
		{
			name:                 "gateway failure with translation failures",
			gatewayFailure:       true,
			translationFailures:  true,
			expectedConfigStatus: clients.ConfigStatusApplyFailed,
		},
		{
			name:                 "both gateway and konnect failure with translation failures",
			gatewayFailure:       true,
			konnectFailure:       true,
			translationFailures:  true,
			expectedConfigStatus: clients.ConfigStatusApplyFailedKonnectApplyFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := clients.CalculateConfigStatus(clients.GatewayConfigApplyStatus{
				ApplyConfigFailed:           tc.gatewayFailure,
				TranslationFailuresOccurred: tc.translationFailures,
			}, clients.KonnectConfigUploadStatus{
				Failed: tc.konnectFailure,
			},
			)
			require.Equal(t, tc.expectedConfigStatus, result)
		})
	}
}
