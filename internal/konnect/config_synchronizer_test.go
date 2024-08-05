package konnect

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

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

// ----------------------------------------------------------------------------
// Mocks: Mock interfaces to run sendconfig.PerformUpdate in tests.
// TODO: These are (mostly) copied from internal/dataplane package, but there are also differences because we do not need so man features for the tests here.
// Should we extract the mock interfaces for sendconfig.PerformUpdate to common packages?
// ----------------------------------------------------------------------------

// mockGatewayClientsProvider is a mock implementation of dataplane.AdminAPIClientsProvider.
type mockGatewayClientsProvider struct {
	gatewayClients []*adminapi.Client
	konnectClient  *adminapi.KonnectClient
	dbMode         dpconf.DBMode
}

func (p *mockGatewayClientsProvider) KonnectClient() *adminapi.KonnectClient {
	return p.konnectClient
}

func (p *mockGatewayClientsProvider) GatewayClients() []*adminapi.Client {
	return p.gatewayClients
}

func (p *mockGatewayClientsProvider) GatewayClientsToConfigure() []*adminapi.Client {
	if p.dbMode.IsDBLessMode() {
		return p.gatewayClients
	}
	if len(p.gatewayClients) == 0 {
		return []*adminapi.Client{}
	}
	return p.gatewayClients[:1]
}

// mockUpdateStrategy is a mock implementation of sendconfig.UpdateStrategyResolver.
type mockUpdateStrategyResolver struct {
	updateCalledForURLs       []string
	lastUpdatedContentForURLs map[string]sendconfig.ContentWithHash
	lock                      sync.RWMutex
}

func newMockUpdateStrategyResolver() *mockUpdateStrategyResolver {
	return &mockUpdateStrategyResolver{
		lastUpdatedContentForURLs: map[string]sendconfig.ContentWithHash{},
	}
}

func (f *mockUpdateStrategyResolver) ResolveUpdateStrategy(
	c sendconfig.UpdateClient,
	_ *diagnostics.ClientDiagnostic,
) sendconfig.UpdateStrategy {
	f.lock.Lock()
	defer f.lock.Unlock()

	url := c.AdminAPIClient().BaseRootURL()
	return &mockUpdateStrategy{onUpdate: f.updateCalledForURLCallback(url)}
}

// updateCalledForURLCallback returns a function that will be called when the mockUpdateStrategy is called.
// That enables us to track which URLs were called.
func (f *mockUpdateStrategyResolver) updateCalledForURLCallback(url string) func(sendconfig.ContentWithHash) error {
	return func(content sendconfig.ContentWithHash) error {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.updateCalledForURLs = append(f.updateCalledForURLs, url)
		f.lastUpdatedContentForURLs[url] = content
		return nil
	}
}

// getUpdateCalledForURLs returns the called URLs.
func (f *mockUpdateStrategyResolver) getUpdateCalledForURLs() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	urls := make([]string, 0, len(f.updateCalledForURLs))
	urls = append(urls, f.updateCalledForURLs...)
	return urls
}

func (f *mockUpdateStrategyResolver) lastUpdatedContentForURL(url string) (sendconfig.ContentWithHash, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	c, ok := f.lastUpdatedContentForURLs[url]
	return c, ok
}

// mockUpdateStrategy is a mock implementation of sendconfig.UpdateStrategy.
type mockUpdateStrategy struct {
	onUpdate func(content sendconfig.ContentWithHash) error
}

func (m *mockUpdateStrategy) Update(_ context.Context, targetContent sendconfig.ContentWithHash) (err error) {
	return m.onUpdate(targetContent)
}

func (m *mockUpdateStrategy) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

func (m *mockUpdateStrategy) Type() string {
	return "Mock"
}

// mockConfigurationChangeDetector is a mock implementation of sendconfig.ConfigurationChangeDetector.
type mockConfigurationChangeDetector struct{}

func (m mockConfigurationChangeDetector) HasConfigurationChanged(
	_ context.Context, oldSHA []byte, newSHA []byte, _ *file.Content, _ sendconfig.KonnectAwareClient, _ sendconfig.StatusClient,
) (bool, error) {
	return !bytes.Equal(oldSHA, newSHA), nil
}

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()

	c, err := adminapi.NewKongAPIClient(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString()), &http.Client{})
	require.NoError(t, err)

	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID, false)
}

// ----------------------------------------------------------------------------
// End of Mocks
// ----------------------------------------------------------------------------

func TestConfigSynchronizer_RunKonnectUpdateServer(t *testing.T) {
	sendConfigPeriod := 10 * time.Millisecond
	testKonnectClient := mustSampleKonnectClient(t)
	resolver := newMockUpdateStrategyResolver()

	s := &ConfigSynchronizer{
		logger:     logr.Discard(),
		syncTicker: time.NewTicker(sendConfigPeriod),
		clientsProvider: &mockGatewayClientsProvider{
			konnectClient: testKonnectClient,
		},
		prometheusMetrics:      metrics.NewCtrlFuncMetrics(),
		updateStrategyResolver: resolver,
		configChangeDetector:   mockConfigurationChangeDetector{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	err := s.Start(ctx)
	require.NoError(t, err)

	t.Logf("Verifying that no URL are updated when no configuration received")
	require.Never(t, func() bool {
		return len(resolver.getUpdateCalledForURLs()) != 0
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
	require.Eventually(t, func() bool {
		urls := resolver.getUpdateCalledForURLs()
		if len(urls) != 1 {
			return false
		}
		url := urls[0]
		contentWithHash, ok := resolver.lastUpdatedContentForURL(url)
		if !ok {
			return false
		}
		return assert.ObjectsAreEqual(content, contentWithHash.Content)
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should send expected configuration in time after received configuration")

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
		urls := resolver.getUpdateCalledForURLs()
		l := len(urls)
		if l == 0 {
			return false
		}
		url := urls[l-1]
		if url != testKonnectClient.BaseRootURL() {
			return false
		}
		contentWithHash, ok := resolver.lastUpdatedContentForURL(url)
		if !ok {
			return false
		}
		return !(assert.ObjectsAreEqual(content, contentWithHash.Content))
	}, 10*sendConfigPeriod, sendConfigPeriod, "Should not send new updates after context cancelled")
}
