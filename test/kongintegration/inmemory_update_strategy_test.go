package kongintegration

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
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

	ctx := t.Context()

	kongC := containers.NewKong(ctx, t)
	kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), managercfg.AdminAPIClientConfig{}, "")
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
	expectedMessage := "invalid path: value must be null"
	// We don't hardcode the expected response body as it can vary between Kong versions
	// Instead, we verify the essential structure and fields are present

	require.EventuallyWithT(t, func(t *assert.CollectT) {
		configSize, err := sut.Update(ctx, faultyConfig)
		if !assert.Error(t, err) {
			return
		}
		// Default value 0 to discard, since error has been returned.
		if !assert.Zero(t, configSize) {
			return
		}
		var updateError sendconfig.UpdateError
		if !assert.True(t, errors.As(err, &updateError)) {
			return
		}
		if wrappedErr := updateError.Unwrap(); !assert.Error(t, wrappedErr) || !assert.IsType(t, &kong.APIError{}, wrappedErr) {
			return
		}
		if !assert.NotEmpty(t, updateError.ResourceFailures()) {
			return
		}
		if !assert.NotEmpty(t, updateError.RawResponseBody()) {
			return
		}
		resourceErr, found := lo.Find(updateError.ResourceFailures(), func(r failures.ResourceFailure) bool {
			return lo.ContainsBy(r.CausingObjects(), func(obj client.Object) bool {
				return obj.GetName() == "test-service"
			})
		})
		if !assert.Truef(t, found, "expected resource error for test-service, got: %+v", updateError.ResourceFailures()) {
			return
		}
		if !assert.Equal(t, resourceErr.Message(), expectedMessage) {
			return
		}
		if diff := cmp.Diff(expectedCausingObjects, resourceErr.CausingObjects()); !assert.Empty(t, diff) {
			return
		}
		// Verify essential structure of the response body
		actualBody := map[string]any{}
		err = json.Unmarshal(updateError.RawResponseBody(), &actualBody)
		if !assert.NoError(t, err) {
			return
		}

		// Check for required fields in the error response
		if !assert.Contains(t, actualBody, "code") {
			return
		}
		if !assert.Contains(t, actualBody, "name") {
			return
		}
		if !assert.Contains(t, actualBody, "flattened_errors") {
			return
		}

		// Verify the flattened_errors contains our test service
		flattenedErrors, ok := actualBody["flattened_errors"].([]any)
		if !assert.True(t, ok) {
			return
		}
		if !assert.NotEmpty(t, flattenedErrors) {
			return
		}

		// Check that at least one error relates to our test service
		foundServiceError := false
		for _, errorItem := range flattenedErrors {
			if errorMap, ok := errorItem.(map[string]any); ok {
				if entityName, ok := errorMap["entity_name"].(string); ok && entityName == "test-service" {
					foundServiceError = true
					break
				}
			}
		}
		if !assert.True(t, foundServiceError, "Expected to find flattened error for test-service") {
			return
		}
	}, timeout, period)
}
