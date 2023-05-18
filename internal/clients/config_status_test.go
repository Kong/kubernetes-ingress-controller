package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
)

func TestChannelConfigNotifier(t *testing.T) {
	logger := testr.New(t)
	n := clients.NewChannelConfigNotifier(logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := n.SubscribeConfigStatus()

	// Call NotifyConfigStatus 5 times to make sure that the method is non-blocking.
	for i := 0; i < 5; i++ {
		n.NotifyConfigStatus(ctx, clients.ConfigStatusOK)
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
			result := clients.CalculateConfigStatus(clients.CalculateConfigStatusInput{
				GatewaysFailed:              tc.gatewayFailure,
				KonnectFailed:               tc.konnectFailure,
				TranslationFailuresOccurred: tc.translationFailures,
			})
			require.Equal(t, tc.expectedConfigStatus, result)
		})
	}
}
