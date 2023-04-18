package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/clients"
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
