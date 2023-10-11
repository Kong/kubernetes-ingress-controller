package kongintegration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
)

// TestUpdateStrategyInMemory_PropagatesResourcesErrors ensures that sendconfig.UpdateStrategyInMemory -
// responsible for executing the configuration update logic - propagates the resources errors returned by the
// Kong Admin API in the flattened errors response.
func TestUpdateStrategyInMemory_PropagatesResourcesErrors(t *testing.T) {
	t.Parallel()

	const (
		timeout = time.Second * 5
		period  = time.Millisecond * 200
	)

	ctx := context.Background()

	adminURL := spawnDBLessKongContainer(ctx, t)
	kongClient, err := kong.NewClient(kong.String(adminURL), &http.Client{})
	require.NoError(t, err)

	sut := sendconfig.NewUpdateStrategyInMemory(
		kongClient,
		sendconfig.DefaultContentToDBLessConfigConverter{},
		logrus.New(),
	)

	// This configuration is faulty and should return a resource error.
	faultyConfig := sendconfig.ContentWithHash{
		Content: &file.Content{
			FormatVersion: "3.0",
			Services: []file.FService{
				{
					Service: kong.Service{
						Name:     kong.String("test-service"),
						Host:     kong.String("konghq.com"),
						Port:     kong.Int(80),
						Protocol: kong.String("grpc"),
						// Paths are not supported for gRPC services. This will trigger an error.
						Path: kong.String("/test"),
						Tags: []*string{
							// Tags are used to identify the resource in the flattened errors response.
							kong.String("k8s-name:test-service"),
							kong.String("k8s-namespace:default"),
							kong.String("k8s-kind:Service"),
							kong.String("k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3"),
							kong.String("k8s-group:"),
							kong.String("k8s-version:v1"),
						},
					},
				},
			},
		},
	}
	expectedResourceError := sendconfig.ResourceError{
		Name:       "test-service",
		Namespace:  "default",
		Kind:       "Service",
		UID:        "a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
		APIVersion: "v1",
		Problems: map[string]string{
			"":     "failed conditional validation given value of field 'protocol'",
			"path": "value must be null",
		},
	}

	require.Eventually(t, func() bool {
		err, resourceErrors, parseErr := sut.Update(ctx, faultyConfig)
		if err == nil {
			t.Logf("expected error: %v", err)
			return false
		}
		if len(resourceErrors) == 0 {
			t.Log("expected resource errors")
			return false
		}
		if parseErr != nil {
			t.Logf("expected no parse error: %v", parseErr)
			return false
		}

		resourceErr, found := lo.Find(resourceErrors, func(r sendconfig.ResourceError) bool {
			return r.Name == "test-service"
		})
		if !found {
			t.Logf("expected resource error for test-service, got: %+v", resourceErrors)
			return false
		}
		if diff := cmp.Diff(expectedResourceError, resourceErr); diff != "" {
			t.Logf("expected resource error to match, got diff: %s", diff)
			return false
		}

		return true
	}, timeout, period)
}
