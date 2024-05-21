package kongintegration

import (
	"context"
	"encoding/json"
	"errors"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
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
	expectedCausingObjects := []client.Object{
		&metav1.PartialObjectMetadata{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-service",
				UID:       "a3b8afcc-9f19-42e4-aa8f-5866168c2ad3",
			},
		},
	}
	expectedMessage := "invalid service:test-service: failed conditional validation given value of field 'protocol'"
	expectedRawErrBody := []byte(`{"code":14,"name":"invalid declarative configuration","fields":{},"message":"declarative config is invalid: {}","flattened_errors":[{"entity_type":"service","entity_name":"test-service","entity_tags":["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"],"errors":[{"type":"field","message":"value must be null","field":"path"},{"type":"entity","message":"failed conditional validation given value of field 'protocol'"}],"entity":{"path":"/test","name":"test-service","protocol":"grpc","tags":["k8s-name:test-service","k8s-namespace:default","k8s-kind:Service","k8s-uid:a3b8afcc-9f19-42e4-aa8f-5866168c2ad3","k8s-group:","k8s-version:v1"],"host":"konghq.com","port":80}}]}`)
	expectedBody := map[string]interface{}{}
	require.NoError(t, json.Unmarshal(expectedRawErrBody, &expectedBody))

	require.Eventually(t, func() bool {
		err := sut.Update(ctx, faultyConfig)
		if err == nil {
			t.Logf("expected error: %v", err)
			return false
		}
		var updateError sendconfig.UpdateError
		if !errors.As(err, &updateError) {
			t.Logf("expected UpdateError, got: %T", err)
			return false
		}
		if len(updateError.ResourceFailures()) == 0 {
			t.Log("expected resource errors")
			return false
		}
		if len(updateError.RawResponseBody()) == 0 {
			t.Log("expected error response body")
			return false
		}
		resourceErr, found := lo.Find(updateError.ResourceFailures(), func(r failures.ResourceFailure) bool {
			return lo.ContainsBy(r.CausingObjects(), func(obj client.Object) bool {
				return obj.GetName() == "test-service"
			})
		})
		if !found {
			t.Logf("expected resource error for test-service, got: %+v", updateError.ResourceFailures())
			return false
		}
		if resourceErr.Message() != expectedMessage {
			t.Logf("expected resource error message to be %q, got %q", expectedMessage, resourceErr.Message())
			return false
		}
		if diff := cmp.Diff(expectedCausingObjects, resourceErr.CausingObjects()); diff != "" {
			t.Logf("expected resource error to match, got diff: %s", diff)
			return false
		}
		actualBody := map[string]interface{}{}
		if err := json.Unmarshal(updateError.RawResponseBody(), &actualBody); err != nil {
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
