package konnect

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
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongAPIClient(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString()), &http.Client{})
	require.NoError(t, err)

	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID, false)
}

func TestConfigSynchronizer_GetTargetContentCopy(t *testing.T) {
	content := &file.Content{
		FormatVersion: "3.0",
		Services: []file.FService{
			{
				Service: kong.Service{
					Name: kong.String("service1"),
					Host: kong.String("example.com"),
				},
				Routes: []*file.FRoute{
					{
						Route: kong.Route{
							Name:       kong.String("route1"),
							Expression: kong.String("http.path == \"/foo\""),
						},
					},
				},
			},
		},
	}

	s := &ConfigSynchronizer{}
	s.SetTargetContent(content)
	copiedContent := s.GetTargetContentCopy()
	require.Equal(t, content, copiedContent, "Copied content should have values with same fields with original content")
	require.NotSame(t, content, copiedContent, "Copied content should not point to the same object with the original content")
}

func TestConfigSynchronizer_RunKonnectUpdateServer(t *testing.T) {
	sendConfigPeriod := 10 * time.Millisecond
	testKonnectClient := mustSampleKonnectClient(t)
	resolver := mocks.NewUpdateStrategyResolver()
	log := logr.Discard()
	s := &ConfigSynchronizer{
		logger:                 logr.Discard(),
		syncTicker:             time.NewTicker(sendConfigPeriod),
		konnectClient:          testKonnectClient,
		prometheusMetrics:      metrics.NewCtrlFuncMetrics(),
		updateStrategyResolver: resolver,
		configChangeDetector:   sendconfig.NewDefaultConfigurationChangeDetector(log),
		configStatusNotifier:   clients.NoOpConfigStatusNotifier{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	err := s.Start(ctx)
	require.NoError(t, err)

	t.Logf("Verifying that no URL are updated when no configuration received")
	require.Never(t, func() bool {
		return len(resolver.GetUpdateCalledForURLs()) != 0
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should not update any URL when no configuration received")

	t.Logf("Verifying that the new config updated when received")
	content := &file.Content{
		FormatVersion: "3.0",
		Services: []file.FService{
			{
				Service: kong.Service{
					Name: kong.String("service1"),
					Host: kong.String("example.com"),
				},
				Routes: []*file.FRoute{
					{
						Route: kong.Route{
							Name:       kong.String("route1"),
							Expression: kong.String("http.path == \"/foo\""),
						},
					},
				},
			},
		},
	}
	s.SetTargetContent(content)
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		urls := resolver.GetUpdateCalledForURLs()
		require.Len(t, urls, 1, "should update only one URL (Konnect)")
		url := urls[0]
		contentWithHash, ok := resolver.LastUpdatedContentForURL(url)
		require.True(t, ok, "should have last updated content for the URL")
		require.Empty(t, cmp.Diff(content, contentWithHash.Content), "should send expected configuration")
	}, 10*sendConfigPeriod, sendConfigPeriod)

	t.Logf("Verifying that update is not called when config not changed")
	l := len(resolver.GetUpdateCalledForURLs())
	s.SetTargetContent(content)
	require.Never(t, func() bool {
		return len(resolver.GetUpdateCalledForURLs()) != l
	}, 10*sendConfigPeriod, sendConfigPeriod)

	t.Logf("Verifying that new config are not sent after context cancelled")
	cancel()
	<-ctx.Done()
	// modify content
	newContent := content.DeepCopy()
	newContent.Consumers = []file.FConsumer{
		{
			Consumer: kong.Consumer{
				Username: kong.String("consumer-1"),
			},
		},
	}
	s.SetTargetContent(newContent)
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
		return !(assert.ObjectsAreEqual(content, contentWithHash.Content))
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should not send new updates after context cancelled")
}
