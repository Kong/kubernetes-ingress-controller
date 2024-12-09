package konnect

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongAPIClient(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString()), &http.Client{})
	require.NoError(t, err)

	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID, false)
}

func TestConfigSynchronizer_RunKonnectUpdateServer(t *testing.T) {
	sendConfigPeriod := 10 * time.Millisecond
	testKonnectClient := mustSampleKonnectClient(t)
	resolver := mocks.NewUpdateStrategyResolver()
	kongConfig := sendconfig.Config{}
	logger := testr.New(t)
	s := NewConfigSynchronizer(
		logger,
		kongConfig,
		sendConfigPeriod,
		&mocks.KonnectClientFactory{Client: testKonnectClient},
		resolver,
		mocks.ConfigurationChangeDetector{},
		clients.NoOpConfigStatusNotifier{},
		mocks.MetricsRecorder{},
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := s.Start(ctx)
		require.NoError(t, err)
	}()

	t.Logf("Verifying that no URL are updated when no configuration received")
	require.Never(t, func() bool {
		return len(resolver.GetUpdateCalledForURLs()) != 0
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should not update any URL when no configuration received")

	t.Logf("Verifying that the new config updated when received")
	expectedContent := &file.Content{
		FormatVersion: "3.0",
		Services: []file.FService{
			{
				Service: kong.Service{
					Name: kong.String("service1"),
					Host: kong.String("example.com"),
				},
			},
		},
	}
	kongState := &kongstate.KongState{
		Services: []kongstate.Service{
			{
				Service: kong.Service{
					Name: kong.String("service1"),
					Host: kong.String("example.com"),
				},
			},
		},
	}
	s.UpdateKongState(ctx, kongState, false)
	// TODO(czeslavo): WHY DOES IT FAIL?
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		urls := resolver.GetUpdateCalledForURLs()
		require.Len(t, urls, 1)

		url := urls[0]
		contentWithHash, ok := resolver.LastUpdatedContentForURL(url)
		require.True(t, ok)
		assert.ObjectsAreEqual(expectedContent, contentWithHash.Content)
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should send expected configuration in time after received configuration")

	t.Logf("Verifying that update is not called when config not changed")
	l := len(resolver.GetUpdateCalledForURLs())
	s.UpdateKongState(ctx, kongState, false)
	require.Never(t, func() bool {
		return len(resolver.GetUpdateCalledForURLs()) != l
	}, 10*sendConfigPeriod, sendConfigPeriod)

	t.Logf("Verifying that new config are not sent after context cancelled")
	cancel()
	<-ctx.Done()
	// modify kong state
	kongState.Services[0].Host = kong.String("example.org")
	expectedContent.Services[0].Host = kong.String("example.org")

	s.UpdateKongState(ctx, kongState, false)
	// The latest updated content should always be the content in the previous update
	// because it should not update new expectedContent after context cancelled.
	require.Never(t, func() bool {
		urls := resolver.GetUpdateCalledForURLs()
		l := len(urls)
		if l == 0 {
			return false
		}
		url := urls[l-1]
		if url != testKonnectClient.BaseRootURL() {
			return false
		}
		contentWithHash, ok := resolver.LastUpdatedContentForURL(url)
		if !ok {
			return false
		}
		return assert.ObjectsAreEqual(expectedContent, contentWithHash.Content)
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should not send new updates after context cancelled")
}
