package kongintegration

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/zapr"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/kongintegration/containers"
)

func TestUpdateStrategyDBMode(t *testing.T) {
	t.Parallel()

	const (
		timeout = time.Second * 5
		period  = time.Millisecond * 200
	)
	ctx := context.Background()

	// Create a network for Postgres and Kong containers to communicate over.
	net, err := network.New(ctx)
	require.NoError(t, err)

	_ = containers.NewPostgres(ctx, t, net)
	kongC := containers.NewKong(ctx, t, containers.KongWithDBMode(net.Name))

	kongClient, err := adminapi.NewKongAPIClient(kongC.AdminURL(ctx, t), &http.Client{})
	require.NoError(t, err)

	logbase, err := zap.NewDevelopment()
	require.NoError(t, err)
	logger := zapr.NewLogger(logbase)
	sut := sendconfig.NewUpdateStrategyDBMode(
		kongClient,
		dump.Config{},
		semver.MustParse("3.6.0"),
		10,
		logger,
	)

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

	const expectedMessage = `invalid service:test-service: HTTP status 400 (message: "2 schema violations (failed conditional validation given value of field 'protocol'; path: value must be null)")`
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		err := sut.Update(ctx, faultyConfig)
		if !assert.Error(t, err) {
			return
		}
		var updateError sendconfig.UpdateError
		if !assert.True(t, errors.As(err, &updateError)) {
			return
		}
		if !assert.NotEmpty(t, updateError.ResourceFailures()) {
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
	}, timeout, period)
	require.NoError(t, err)
}
