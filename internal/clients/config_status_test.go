package clients

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
)

func TestChannelConfigNotifier(t *testing.T) {
	logger := testr.New(t)
	n := NewChannelConfigNotifier(logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := n.SubscribeConfigStatus()

	// Call NotifyConfigStatus 5 times to make sure that the method is non-blocking.
	for i := 0; i < 5; i++ {
		n.NotifyConfigStatus(ctx, ConfigStatusOK)
	}

	for i := 0; i < 5; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for config status i=%d", i)
		}
	}
}
