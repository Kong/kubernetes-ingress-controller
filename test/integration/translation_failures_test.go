//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestTranslationFailures ensures that proper warning Kubernetes events are recorded in case of translation failures
// encountered.
func TestTranslationFailures(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	testCases := []struct {
		name                      string
		translationFailureTrigger func()
		expectEventsForObjects    []string
	}{}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.translationFailureTrigger()

			require.Eventually(t, func() bool {
				eventsForAllObjectsFound := true
				for _, objectName := range tt.expectEventsForObjects {
					events, err := env.Cluster().Client().CoreV1().Events(ns.GetName()).List(ctx, metav1.ListOptions{
						FieldSelector: fmt.Sprintf("reason=%s", dataplane.KongConfigurationTranslationFailedEventReason),
					})
					if err != nil {
						t.Logf("failed to list events: %s", err)
						eventsForAllObjectsFound = false
					}

					if len(events.Items) == 0 {
						t.Logf("waiting for events related to %s to be created", objectName)
						eventsForAllObjectsFound = false
					}
				}
				return eventsForAllObjectsFound
			}, time.Minute, time.Second)
		})
	}
}
