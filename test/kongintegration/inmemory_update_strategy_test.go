package kongintegration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
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

	kongC := containers.NewKong(ctx, t)
	kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), &http.Client{})
	require.NoError(t, err)

	logbase, err := zap.NewDevelopment()
	require.NoError(t, err)
	logger := zapr.NewLogger(logbase)
	sut := sendconfig.NewUpdateStrategyInMemory(
		kongClient,
		sendconfig.DefaultContentToDBLessConfigConverter{},
		logger,
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
			"service:test-service": "failed conditional validation given value of field 'protocol'",
			"path":                 "value must be null",
		},
	}
	expectedRawErrBody := []byte(`{"code":14,"name":"invalid declarative configuration","fields":{},"message":"declarative config is invalid: {}","flattened_errors":[{"entity_type":"service","entity_name":"test-service","entity_tags":["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"],"errors":[{"type":"field","message":"value must be null","field":"path"},{"type":"entity","message":"failed conditional validation given value of field 'protocol'"}],"entity":{"path":"/test","name":"test-service","protocol":"grpc","tags":["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"],"host":"konghq.com","port":80}}]}`)
	expectedBody := map[string]interface{}{}
	require.NoError(t, json.Unmarshal(expectedRawErrBody, &expectedBody))

	require.Eventually(t, func() bool {
		err, resourceErrors, rawErrBody, parseErr := sut.Update(ctx, faultyConfig)
		if err == nil {
			t.Logf("expected error: %v", err)
			return false
		}
		if len(resourceErrors) == 0 {
			t.Log("expected resource errors")
			return false
		}
		if len(rawErrBody) == 0 {
			t.Log("expected error response body")
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
		actualBody := map[string]interface{}{}
		if err := json.Unmarshal(rawErrBody, &actualBody); err != nil {
			t.Logf("could not unmarshal error body: %s", err)
			return false
		}
		if diff := cmp.Diff(expectedBody, actualBody); diff != "" {
			t.Logf("unexpected error body diff: %s", diff)
			return false
		}

		return true
	}, timeout, period)
}
