package dataplane

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/test/mocks"
)

var defaultKongStatus = kong.Status{
	ConfigurationHash: sendconfig.WellKnownInitialHash,
}

func TestUniqueObjects(t *testing.T) {
	t.Log("generating some objects to test the de-duplication of objects")
	ing1 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  corev1.NamespaceDefault,
			Name:       "test-ingress-1",
			Generation: 1,
		},
	}
	ing1.SetGroupVersionKind(ingGVK)
	ing2 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  corev1.NamespaceDefault,
			Name:       "test-ingress-2",
			Generation: 1,
		},
	}
	ing2.SetGroupVersionKind(ingGVK)
	ing3 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  "other-namespace",
			Name:       "test-ingress-1",
			Generation: 1,
		},
	}
	ing3.SetGroupVersionKind(ingGVK)
	ing4 := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  "other-namespace",
			Name:       "test-ingress-2",
			Generation: 1,
		},
	}
	ing4.SetGroupVersionKind(ingGVK)

	testCases := []struct {
		name         string
		reportedObjs []client.Object
		failedObjs   [][]client.Object
		uniqueObjs   []client.Object
	}{
		{
			name:         "no failures",
			reportedObjs: []client.Object{ing1, ing2},
			uniqueObjs:   []client.Object{ing1, ing2},
		},
		{
			name:         "has failures",
			reportedObjs: []client.Object{ing1, ing3},
			failedObjs: [][]client.Object{
				{ing1},
				{ing4},
			},
			uniqueObjs: []client.Object{ing1, ing3, ing4},
		},
		{
			name:         "one object in multiple failures",
			reportedObjs: []client.Object{ing1, ing2},
			failedObjs: [][]client.Object{
				{ing3},
				{ing2, ing3},
			},
			uniqueObjs: []client.Object{ing1, ing2, ing3},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			translationFailures := []failures.ResourceFailure{}
			for _, failedObjs := range tc.failedObjs {
				translationFailure, err := failures.NewResourceFailure(
					"for test", failedObjs...,
				)
				require.NoError(t, err)
				translationFailures = append(translationFailures, translationFailure)
			}
			uniqueObjs := UniqueObjects(tc.reportedObjs, translationFailures)
			require.Len(t, uniqueObjs, len(tc.uniqueObjs))
			require.ElementsMatch(t, tc.uniqueObjs, uniqueObjs)
		})
	}
}

// initialized objects don't have GVK's, so we fake those for unit tests.
var (
	ingGVK = schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	}
)

// mockGatewayClientsProvider is a mock implementation of dataplane.AdminAPIClientsProvider.
type mockGatewayClientsProvider struct {
	gatewayClients []*adminapi.Client
	konnectClient  *adminapi.KonnectClient
}

func (f mockGatewayClientsProvider) KonnectClient() *adminapi.KonnectClient {
	return f.konnectClient
}

func (f mockGatewayClientsProvider) GatewayClients() []*adminapi.Client {
	return f.gatewayClients
}

// mockUpdateStrategy is a mock implementation of sendconfig.UpdateStrategyResolver.
type mockUpdateStrategyResolver struct {
	updateCalledForURLs       []string
	lastUpdatedContentForURLs map[string]sendconfig.ContentWithHash
	shouldReturnErrorOnUpdate map[string]struct{}
	t                         *testing.T
	lock                      sync.RWMutex
	singleError               bool
}

func newMockUpdateStrategyResolver(t *testing.T) *mockUpdateStrategyResolver {
	return &mockUpdateStrategyResolver{
		t:                         t,
		shouldReturnErrorOnUpdate: map[string]struct{}{},
		lastUpdatedContentForURLs: map[string]sendconfig.ContentWithHash{},
	}
}

func (f *mockUpdateStrategyResolver) ResolveUpdateStrategy(c sendconfig.UpdateClient) sendconfig.UpdateStrategy {
	f.lock.Lock()
	defer f.lock.Unlock()

	url := c.AdminAPIClient().BaseRootURL()
	return &mockUpdateStrategy{onUpdate: f.updateCalledForURLCallback(url, f.singleError)}
}

// returnErrorOnUpdate will cause the mockUpdateStrategy with a given Admin API URL to return an error on Update().
func (f *mockUpdateStrategyResolver) returnErrorOnUpdate(url string, shouldReturnErr bool) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if shouldReturnErr {
		f.shouldReturnErrorOnUpdate[url] = struct{}{}
	} else {
		delete(f.shouldReturnErrorOnUpdate, url)
	}
}

// updateCalledForURLCallback returns a function that will be called when the mockUpdateStrategy is called.
// That enables us to track which URLs were called.
func (f *mockUpdateStrategyResolver) updateCalledForURLCallback(url string, singleError bool) func(sendconfig.ContentWithHash) error {
	return func(content sendconfig.ContentWithHash) error {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.updateCalledForURLs = append(f.updateCalledForURLs, url)
		if _, ok := f.shouldReturnErrorOnUpdate[url]; ok {
			if singleError {
				delete(f.shouldReturnErrorOnUpdate, url)
			}
			return errors.New("error on update")
		}
		f.lastUpdatedContentForURLs[url] = content

		return nil
	}
}

// assertUpdateCalledForURLs asserts that the mockUpdateStrategy was called for the given URLs.
func (f *mockUpdateStrategyResolver) assertUpdateCalledForURLs(urls []string) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	require.ElementsMatch(f.t, urls, f.updateCalledForURLs, "update was not called for all URLs")
}

func (f *mockUpdateStrategyResolver) assertNoUpdateCalled() {
	f.lock.RLock()
	defer f.lock.RUnlock()

	require.Empty(f.t, f.updateCalledForURLs, "update was called")
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

func (m *mockUpdateStrategy) Update(_ context.Context, content sendconfig.ContentWithHash) (
	err error,
	resourceErrors []sendconfig.ResourceError,
	resourceErrorsParseErr error,
) {
	err = m.onUpdate(content)
	return err, nil, nil
}

func (m *mockUpdateStrategy) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDBLess
}

func (m *mockUpdateStrategy) Type() string {
	return "Mock"
}

// mockConfigurationChangeDetector is a mock implementation of sendconfig.ConfigurationChangeDetector.
type mockConfigurationChangeDetector struct {
	hasConfigurationChanged bool
	status                  kong.Status
}

func (m mockConfigurationChangeDetector) HasConfigurationChanged(
	context.Context, []byte, []byte, *file.Content, sendconfig.KonnectAwareClient, sendconfig.StatusClient,
) (bool, error) {
	return m.hasConfigurationChanged, nil
}

func TestKongClientUpdate_AllExpectedClientsAreCalledAndErrorIsPropagated(t *testing.T) {
	var (
		ctx                = context.Background()
		testKonnectClient  = mustSampleKonnectClient(t)
		testGatewayClients = []*adminapi.Client{
			mustSampleGatewayClient(t),
			mustSampleGatewayClient(t),
		}
	)

	testCases := []struct {
		name                 string
		gatewayClients       []*adminapi.Client
		konnectClient        *adminapi.KonnectClient
		errorOnUpdateForURLs []string
		expectError          bool
	}{
		{
			name:           "2 gateway clients and konnect with no errors",
			gatewayClients: testGatewayClients,
			konnectClient:  testKonnectClient,
			expectError:    false,
		},
		{
			name:                 "2 gateway clients and konnect with error on konnect",
			gatewayClients:       testGatewayClients,
			konnectClient:        testKonnectClient,
			errorOnUpdateForURLs: []string{testKonnectClient.BaseRootURL()},
			expectError:          false,
		},
		{
			name:                 "2 gateway clients with error on one of them",
			gatewayClients:       testGatewayClients,
			errorOnUpdateForURLs: []string{testGatewayClients[0].BaseRootURL()},
			expectError:          true,
		},
		{
			name:           "2 gateway clients and konnect with error on one of gateways and konnect",
			gatewayClients: testGatewayClients,
			errorOnUpdateForURLs: []string{
				testGatewayClients[0].BaseRootURL(),
				testKonnectClient.BaseRootURL(),
			},
			expectError: true,
		},
		{
			name:          "only konnect client with no error",
			konnectClient: testKonnectClient,
			expectError:   false,
		},
		{
			name:                 "only konnect client with error on it",
			konnectClient:        testKonnectClient,
			errorOnUpdateForURLs: []string{testKonnectClient.BaseRootURL()},
			expectError:          false,
		},
		{
			name:        "no clients at all",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			clientsProvider := mockGatewayClientsProvider{
				gatewayClients: tc.gatewayClients,
				konnectClient:  tc.konnectClient,
			}
			updateStrategyResolver := newMockUpdateStrategyResolver(t)
			for _, url := range tc.errorOnUpdateForURLs {
				updateStrategyResolver.returnErrorOnUpdate(url, true)
			}
			// always return true for HasConfigurationChanged to trigger an update
			configChangeDetector := mockConfigurationChangeDetector{
				hasConfigurationChanged: true,
				status:                  defaultKongStatus,
			}
			configBuilder := newMockKongConfigBuilder()
			kongRawStateGetter := &mockKongLastValidConfigFetcher{}
			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)

			err := kongClient.Update(ctx)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			allExpectedURLs := mapClientsToUrls(clientsProvider)
			updateStrategyResolver.assertUpdateCalledForURLs(allExpectedURLs)
		})
	}
}

func TestKongClientUpdate_WhenNoChangeInConfigNoClientGetsCalled(t *testing.T) {
	clientsProvider := mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{
			mustSampleGatewayClient(t),
			mustSampleGatewayClient(t),
		},
		konnectClient: mustSampleKonnectClient(t),
	}
	updateStrategyResolver := newMockUpdateStrategyResolver(t)

	// no change in config, we'll expect no update to be called
	configChangeDetector := mockConfigurationChangeDetector{
		hasConfigurationChanged: false,
		status:                  defaultKongStatus,
	}
	configBuilder := newMockKongConfigBuilder()
	kongRawStateGetter := &mockKongLastValidConfigFetcher{}
	kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)

	ctx := context.Background()
	err := kongClient.Update(ctx)
	require.NoError(t, err)

	updateStrategyResolver.assertNoUpdateCalled()
}

type mockConfigStatusQueue struct {
	notifications []clients.ConfigStatus
	lock          sync.RWMutex
}

func newMockConfigStatusQueue() *mockConfigStatusQueue {
	return &mockConfigStatusQueue{}
}

func (m *mockConfigStatusQueue) NotifyConfigStatus(_ context.Context, status clients.ConfigStatus) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.notifications = append(m.notifications, status)
}

func (m *mockConfigStatusQueue) Notifications() []clients.ConfigStatus {
	m.lock.RLock()
	defer m.lock.RUnlock()
	copied := make([]clients.ConfigStatus, len(m.notifications))
	copy(copied, m.notifications)
	return copied
}

type mockKongConfigBuilder struct {
	translationFailuresToReturn []failures.ResourceFailure
	kongState                   *kongstate.KongState
}

func newMockKongConfigBuilder() *mockKongConfigBuilder {
	return &mockKongConfigBuilder{
		kongState: &kongstate.KongState{},
	}
}

func (p *mockKongConfigBuilder) BuildKongConfig() parser.KongConfigBuildingResult {
	return parser.KongConfigBuildingResult{
		KongState:           p.kongState,
		TranslationFailures: p.translationFailuresToReturn,
	}
}

func (p *mockKongConfigBuilder) returnTranslationFailures(enabled bool) {
	if enabled {
		// Return some mocked translation failures.
		p.translationFailuresToReturn = []failures.ResourceFailure{
			lo.Must(failures.NewResourceFailure("some reason", &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "name",
					Namespace: "namespace",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
			},
			)),
		}
	} else {
		p.translationFailuresToReturn = nil
	}
}

func TestKongClientUpdate_ConfigStatusIsNotified(t *testing.T) {
	var (
		ctx               = context.Background()
		testKonnectClient = mustSampleKonnectClient(t)
		testGatewayClient = mustSampleGatewayClient(t)

		clientsProvider = mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{testGatewayClient},
			konnectClient:  testKonnectClient,
		}

		updateStrategyResolver = newMockUpdateStrategyResolver(t)
		configChangeDetector   = mockConfigurationChangeDetector{hasConfigurationChanged: true}
		configBuilder          = newMockKongConfigBuilder()
		kongRawStateGetter     = &mockKongLastValidConfigFetcher{}
		kongClient             = setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)
	)

	testCases := []struct {
		name                string
		gatewayFailure      bool
		konnectFailure      bool
		translationFailures bool
		expectedStatus      clients.ConfigStatus
	}{
		{
			name:                "success",
			gatewayFailure:      false,
			konnectFailure:      false,
			translationFailures: false,
			expectedStatus:      clients.ConfigStatusOK,
		},
		{
			name:                "gateway failure",
			gatewayFailure:      true,
			konnectFailure:      false,
			translationFailures: false,
			expectedStatus:      clients.ConfigStatusApplyFailed,
		},
		{
			name:                "translation failures",
			gatewayFailure:      false,
			konnectFailure:      false,
			translationFailures: true,
			expectedStatus:      clients.ConfigStatusTranslationErrorHappened,
		},
		{
			name:                "konnect failure",
			gatewayFailure:      false,
			konnectFailure:      true,
			translationFailures: false,
			expectedStatus:      clients.ConfigStatusOKKonnectApplyFailed,
		},
		{
			name:                "both gateway and konnect failure",
			gatewayFailure:      true,
			konnectFailure:      true,
			translationFailures: false,
			expectedStatus:      clients.ConfigStatusApplyFailedKonnectApplyFailed,
		},
		{
			name:                "translation failures and konnect failure",
			gatewayFailure:      false,
			konnectFailure:      true,
			translationFailures: true,
			expectedStatus:      clients.ConfigStatusTranslationErrorHappenedKonnectApplyFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the status queue. We want to make sure that the status is always notified.
			statusQueue := newMockConfigStatusQueue()
			kongClient.SetConfigStatusNotifier(statusQueue)

			updateStrategyResolver.returnErrorOnUpdate(testGatewayClient.BaseRootURL(), tc.gatewayFailure)
			updateStrategyResolver.returnErrorOnUpdate(testKonnectClient.BaseRootURL(), tc.konnectFailure)
			configBuilder.returnTranslationFailures(tc.translationFailures)

			_ = kongClient.Update(ctx)
			notifications := statusQueue.Notifications()
			require.Len(t, notifications, 1)
			require.Equal(t, tc.expectedStatus, notifications[0])

			_ = kongClient.Update(ctx)
			notifications = statusQueue.Notifications()
			require.Len(t, notifications, 1, "no new notification should be sent if the status hasn't changed")
		})
	}
}

func TestKongClient_ApplyConfigurationEvents(t *testing.T) {
	testGatewayClient := mustSampleGatewayClient(t)
	clientsProvider := mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{testGatewayClient},
	}
	updateStrategyResolver := newMockUpdateStrategyResolver(t)
	configChangeDetector := mockConfigurationChangeDetector{hasConfigurationChanged: true}
	configBuilder := newMockKongConfigBuilder()

	testCases := []struct {
		name         string
		podNamespace string
		podName      string
		expectEvents bool
	}{
		{
			name:         "events recorded when POD_NAMESPACE and POD_NAME are set",
			podNamespace: "test-namespace",
			podName:      "test-pod",
			expectEvents: true,
		},
		{
			name:         "no events when POD_NAMESPACE and POD_NAME are not set",
			expectEvents: false,
		},
		{
			name:         "no events when POD_NAMESPACE is not set",
			podName:      "test-pod",
			expectEvents: false,
		},
		{
			name:         "no events when POD_NAME is not set",
			podNamespace: "test-namespace",
			expectEvents: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("POD_NAMESPACE", tc.podNamespace)
			t.Setenv("POD_NAME", tc.podName)

			eventRecorder := mocks.NewEventRecorder()
			kongRawStateGetter := &mockKongLastValidConfigFetcher{}
			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, eventRecorder, kongRawStateGetter)

			err := kongClient.Update(context.Background())
			require.NoError(t, err)

			if tc.expectEvents {
				require.NotEmpty(t, eventRecorder.Events())
			} else {
				require.Empty(t, eventRecorder.Events())
			}
		})
	}
}

func TestKongClient_EmptyConfigUpdate(t *testing.T) {
	var (
		ctx               = context.Background()
		testKonnectClient = mustSampleKonnectClient(t)
		testGatewayClient = mustSampleGatewayClient(t)

		clientsProvider = mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{testGatewayClient},
			konnectClient:  testKonnectClient,
		}

		updateStrategyResolver = newMockUpdateStrategyResolver(t)
		configChangeDetector   = mockConfigurationChangeDetector{hasConfigurationChanged: true}
		configBuilder          = newMockKongConfigBuilder()
		kongRawStateGetter     = &mockKongLastValidConfigFetcher{}
		kongClient             = setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)
	)

	t.Run("dbless", func(t *testing.T) {
		kongClient.kongConfig.InMemory = true
		err := kongClient.Update(ctx)
		require.NoError(t, err)

		gwContent, ok := updateStrategyResolver.lastUpdatedContentForURL(testGatewayClient.BaseRootURL())
		require.True(t, ok)
		assert.Equal(t, gwContent.Content, &file.Content{
			Upstreams: []file.FUpstream{
				{
					Upstream: kong.Upstream{
						Name: lo.ToPtr(deckgen.StubUpstreamName),
					},
				},
			},
		}, "gateway content should have appended stub upstream")

		konnectContent, ok := updateStrategyResolver.lastUpdatedContentForURL(testKonnectClient.BaseRootURL())
		require.True(t, ok)
		assert.Equal(t, konnectContent.Content, &file.Content{}, "konnect content should be empty")
	})

	t.Run("db", func(t *testing.T) {
		kongClient.kongConfig.InMemory = false
		err := kongClient.Update(ctx)
		require.NoError(t, err)

		gwContent, ok := updateStrategyResolver.lastUpdatedContentForURL(testGatewayClient.BaseRootURL())
		require.True(t, ok)
		assert.Equal(t, gwContent.Content, &file.Content{}, "gateway content should be empty")

		konnectContent, ok := updateStrategyResolver.lastUpdatedContentForURL(testKonnectClient.BaseRootURL())
		require.True(t, ok)
		assert.Equal(t, konnectContent.Content, &file.Content{}, "konnect content should be empty")
	})
}

// setupTestKongClient creates a KongClient with mocked dependencies.
func setupTestKongClient(
	t *testing.T,
	updateStrategyResolver *mockUpdateStrategyResolver,
	clientsProvider mockGatewayClientsProvider,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
	configBuilder *mockKongConfigBuilder,
	eventRecorder record.EventRecorder,
	kongRawStateGetter configfetcher.LastValidConfigFetcher,
) *KongClient {
	logger := logrus.New()
	timeout := time.Second
	ingressClass := "kong"
	diagnostic := util.ConfigDumpDiagnostic{}
	config := sendconfig.Config{}
	dbMode := "off"

	if eventRecorder == nil {
		eventRecorder = mocks.NewEventRecorder()
	}

	kongClient, err := NewKongClient(
		logger,
		timeout,
		ingressClass,
		diagnostic,
		config,
		eventRecorder,
		dbMode,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
		kongRawStateGetter,
		configBuilder,
		store.NewCacheStores(),
	)
	require.NoError(t, err)
	return kongClient
}

func mustSampleGatewayClient(t *testing.T) *adminapi.Client {
	t.Helper()
	c, err := adminapi.NewTestClient(fmt.Sprintf("https://%s:8080", uuid.NewString()))
	require.NoError(t, err)
	return c
}

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()

	c, err := kong.NewClient(lo.ToPtr(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString())), &http.Client{})
	require.NoError(t, err)

	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID)
}

func mapClientsToUrls(clients mockGatewayClientsProvider) []string {
	urls := lo.Map(clients.GatewayClients(), func(c *adminapi.Client, _ int) string {
		return c.BaseRootURL()
	})
	if clients.KonnectClient() != nil {
		urls = append(urls, clients.KonnectClient().BaseRootURL())
	}
	return urls
}

type mockKongLastValidConfigFetcher struct {
	kongRawState  *utils.KongRawState
	lastKongState *kongstate.KongState
	status        kong.Status
}

func (cf *mockKongLastValidConfigFetcher) LastValidConfig() (*kongstate.KongState, bool) {
	if cf.lastKongState != nil {
		return cf.lastKongState, true
	}
	return nil, false
}

func (cf *mockKongLastValidConfigFetcher) StoreLastValidConfig(s *kongstate.KongState) {
	cf.lastKongState = s
}

func (cf *mockKongLastValidConfigFetcher) TryFetchingValidConfigFromGateways(context.Context, logrus.FieldLogger, []*adminapi.Client) error {
	if cf.kongRawState != nil {
		cf.lastKongState = configfetcher.KongRawStateToKongState(cf.kongRawState)
	}
	return nil
}

func TestKongClientUpdate_FetchStoreAndPushLastValidConfig(t *testing.T) {
	var (
		ctx = context.Background()

		clientsProvider = mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{
				mustSampleGatewayClient(t),
				mustSampleGatewayClient(t),
			},
		}

		updateStrategyResolver = newMockUpdateStrategyResolver(t)
		configChangeDetector   = mockConfigurationChangeDetector{hasConfigurationChanged: true}
		lastKongRawState       = &utils.KongRawState{
			Services: []*kong.Service{
				{
					Name: kong.String("last_service"),
					ID:   kong.String("abc"),
				},
			},
			Routes: []*kong.Route{
				{
					Name: kong.String("last_route"),
					Service: &kong.Service{
						ID: kong.String("abc"),
					},
				},
			},
		}
		lastKongState = configfetcher.KongRawStateToKongState(lastKongRawState)
		newKongState  = &kongstate.KongState{
			Services: []kongstate.Service{
				{
					Service: kong.Service{
						Name: kong.String("new_service"),
					},
					Namespace: "new_namespace",
					Routes: []kongstate.Route{
						{
							Route: kong.Route{
								Name: kong.String("new_route"),
							},
						},
					},
				},
			},
		}
		configBuilder = newMockKongConfigBuilder()
	)
	configBuilder.kongState = newKongState

	testCases := []struct {
		name                  string
		gatewayFailure        bool
		translationFailures   bool
		singleError           bool
		lastValidKongRawState *utils.KongRawState
		lastKongStatusHash    string
		expectedLastKongState *kongstate.KongState
		errorsSize            int
	}{
		{
			name:                  "success, new fallback set",
			lastValidKongRawState: lastKongRawState,
			expectedLastKongState: newKongState,
			lastKongStatusHash:    "xyz",
		},
		{
			name:                  "no previous state, failure",
			gatewayFailure:        true,
			singleError:           true,
			expectedLastKongState: nil,
			errorsSize:            1,
			lastKongStatusHash:    sendconfig.WellKnownInitialHash,
		},
		{
			name:                  "previous state, failure, fallback pushed with success",
			gatewayFailure:        true,
			singleError:           true,
			lastValidKongRawState: lastKongRawState,
			expectedLastKongState: lastKongState,
			errorsSize:            1,
			lastKongStatusHash:    "xyz",
		},
		{
			name:                  "previous state, failure, fallback pushed with failure",
			gatewayFailure:        true,
			singleError:           false,
			lastValidKongRawState: lastKongRawState,
			expectedLastKongState: lastKongState,
			errorsSize:            3,
			lastKongStatusHash:    "xyz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategyResolver.returnErrorOnUpdate(clientsProvider.gatewayClients[0].BaseRootURL(), tc.gatewayFailure)
			updateStrategyResolver.returnErrorOnUpdate(clientsProvider.gatewayClients[1].BaseRootURL(), tc.gatewayFailure)

			updateStrategyResolver.singleError = tc.singleError
			configChangeDetector.status.ConfigurationHash = tc.lastKongStatusHash
			kongRawStateGetter := &mockKongLastValidConfigFetcher{
				status:       configChangeDetector.status,
				kongRawState: tc.lastValidKongRawState,
			}
			kongClient := setupTestKongClient(
				t,
				updateStrategyResolver,
				clientsProvider,
				configChangeDetector,
				configBuilder,
				nil,
				kongRawStateGetter,
			)

			err := kongClient.Update(ctx)
			if tc.errorsSize > 0 {
				// check if the error is joined with other errors. When there are multiple errors,
				// they are separated by \n, hence we count the number of \n.
				assert.Equal(t, tc.errorsSize, strings.Count(err.Error(), "\n"))
			} else {
				assert.NoError(t, err)
			}
			s, _ := kongClient.kongConfigFetcher.LastValidConfig()
			assert.Equal(t, tc.expectedLastKongState, s)
		})
	}
}
