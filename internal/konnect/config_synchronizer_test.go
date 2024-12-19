package konnect_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/konnect"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()
	c, err := adminapi.NewKongAPIClient(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString()), &http.Client{})
	require.NoError(t, err)
	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID, false)
}

func TestConfigSynchronizer_UpdatesKongConfigAccordingly(t *testing.T) {
	log := logr.Discard()
	sendConfigPeriod := 10 * time.Millisecond
	testKonnectClient := mustSampleKonnectClient(t)
	resolver := mocks.NewUpdateStrategyResolver()
	configStatusNotifier := clients.NewChannelConfigNotifier(log)
	s := konnect.NewConfigSynchronizer(
		konnect.ConfigSynchronizerParams{
			Logger:                 log,
			ConfigUploadPeriod:     sendConfigPeriod,
			KonnectClientFactory:   &mocks.KonnectClientFactory{Client: testKonnectClient},
			UpdateStrategyResolver: resolver,
			ConfigChangeDetector:   sendconfig.NewDefaultConfigurationChangeDetector(log),
			ConfigStatusNotifier:   configStatusNotifier,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	kongState := func() *kongstate.KongState {
		return &kongstate.KongState{
			Services: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("service1"),
						Host: kong.String("example.com"),
					},
				},
			},
		}
	}
	s.UpdateKongState(ctx, kongState(), false)
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		urls := resolver.GetUpdateCalledForURLs()
		require.Len(t, urls, 1, "should update only one URL (Konnect)")
		url := urls[0]
		contentWithHash, ok := resolver.LastUpdatedContentForURL(url)
		require.True(t, ok, "should have last updated content for the URL")
		require.Empty(t, cmp.Diff(expectedContent, contentWithHash.Content), "should send expected configuration")
	}, 10*sendConfigPeriod, sendConfigPeriod)

	t.Logf("Verifying that update is not called when config not changed")
	l := len(resolver.GetUpdateCalledForURLs())
	s.UpdateKongState(ctx, kongState(), false)
	require.Never(t, func() bool {
		return len(resolver.GetUpdateCalledForURLs()) != l
	}, 10*sendConfigPeriod, sendConfigPeriod)

	t.Logf("Verifying that new config are not sent after context cancelled")
	cancel()
	<-ctx.Done()

	// Modify the Kong state and expected content and update it again.
	state := kongState()
	state.Services[0].Host = kong.String("example.org")
	expectedContent.Services[0].Host = kong.String("example.org")
	s.UpdateKongState(ctx, state, false)

	// The latest updated content should always be the content in the previous update
	// because it should not update new content after context cancelled.
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

func TestConfigSynchronizer_ConfigIsSanitizedWhenConfiguredSo(t *testing.T) {
	log := logr.Discard()
	sendConfigPeriod := 10 * time.Millisecond
	testKonnectClient := mustSampleKonnectClient(t)
	resolver := mocks.NewUpdateStrategyResolver()
	configStatusNotifier := clients.NewChannelConfigNotifier(log)
	s := konnect.NewConfigSynchronizer(
		konnect.ConfigSynchronizerParams{
			Logger: log,
			KongConfig: sendconfig.Config{
				SanitizeKonnectConfigDumps: true,
			},
			ConfigUploadPeriod:     sendConfigPeriod,
			KonnectClientFactory:   &mocks.KonnectClientFactory{Client: testKonnectClient},
			UpdateStrategyResolver: resolver,
			ConfigChangeDetector:   sendconfig.NewDefaultConfigurationChangeDetector(log),
			ConfigStatusNotifier:   configStatusNotifier,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := s.Start(ctx)
		require.NoError(t, err)
	}()

	t.Log("Updating Kong state with sensitive information")
	kongState := &kongstate.KongState{
		Certificates: []kongstate.Certificate{
			{
				Certificate: kong.Certificate{
					ID:  kong.String("new_cert"),
					Key: kong.String(`private-key-string`), // This should be redacted.
				},
			},
		},
	}
	s.UpdateKongState(ctx, kongState, false)

	t.Log("Verifying that the sensitive information is redacted")
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		konnectContent, ok := resolver.LastUpdatedContentForURL(testKonnectClient.BaseRootURL())
		require.True(t, ok, "should have last updated content for the URL")

		require.Len(t, konnectContent.Content.Certificates, 1, "expected 1 certificate")
		cert := konnectContent.Content.Certificates[0]
		require.NotNil(t, cert.Key, "expected certificate key")
		require.Equal(t, "{vault://redacted-value}", *cert.Key, "expected redacted certificate key")
	}, 10*sendConfigPeriod, sendConfigPeriod)
}
