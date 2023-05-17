package dataplane_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/deck/file"
	gokong "github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

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
			uniqueObjs := dataplane.UniqueObjects(tc.reportedObjs, translationFailures)
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
	shouldReturnErrorOnUpdate map[string]struct{}
	t                         *testing.T
	lock                      sync.RWMutex
}

func newMockUpdateStrategyResolver(t *testing.T) *mockUpdateStrategyResolver {
	return &mockUpdateStrategyResolver{
		t:                         t,
		shouldReturnErrorOnUpdate: map[string]struct{}{},
	}
}

func (f *mockUpdateStrategyResolver) ResolveUpdateStrategy(c sendconfig.UpdateClient) sendconfig.UpdateStrategy {
	f.lock.Lock()
	defer f.lock.Unlock()

	url := c.AdminAPIClient().BaseRootURL()
	return &mockUpdateStrategy{onUpdate: f.updateCalledForURLCallback(url)}
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
func (f *mockUpdateStrategyResolver) updateCalledForURLCallback(url string) func() error {
	return func() error {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.updateCalledForURLs = append(f.updateCalledForURLs, url)
		if _, ok := f.shouldReturnErrorOnUpdate[url]; ok {
			return errors.New("error on update")
		}
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

// mockUpdateStrategy is a mock implementation of sendconfig.UpdateStrategy.
type mockUpdateStrategy struct {
	onUpdate func() error
}

func (m *mockUpdateStrategy) Update(context.Context, sendconfig.ContentWithHash) (
	err error,
	resourceErrors []sendconfig.ResourceError,
	resourceErrorsParseErr error,
) {
	err = m.onUpdate()
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
}

func (m mockConfigurationChangeDetector) HasConfigurationChanged(
	context.Context, []byte, []byte, *file.Content, sendconfig.KonnectAwareClient, sendconfig.StatusClient,
) (bool, error) {
	return m.hasConfigurationChanged, nil
}

func TestKongClientUpdate_AllExpectedClientsAreCalledAndErrorIsPropagated(t *testing.T) {
	t.Parallel()

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
			configChangeDetector := mockConfigurationChangeDetector{hasConfigurationChanged: true}
			configBuilder := newMockKongConfigBuilder()

			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder)

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
	configChangeDetector := mockConfigurationChangeDetector{hasConfigurationChanged: false}
	configBuilder := newMockKongConfigBuilder()

	kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder)

	ctx := context.Background()
	err := kongClient.Update(ctx)
	require.NoError(t, err)

	updateStrategyResolver.assertNoUpdateCalled()
}

type mockConfigStatusQueue struct {
	wasNotified bool
}

func newMockConfigStatusQueue() *mockConfigStatusQueue {
	return &mockConfigStatusQueue{}
}

func (m *mockConfigStatusQueue) NotifyConfigStatus(context.Context, clients.ConfigStatus) {
	m.wasNotified = true
}

func (m *mockConfigStatusQueue) WasNotified() bool {
	return m.wasNotified
}

type mockKongConfigBuilder struct {
	translationFailuresToReturn []failures.ResourceFailure
}

func newMockKongConfigBuilder() *mockKongConfigBuilder {
	return &mockKongConfigBuilder{}
}

func (p *mockKongConfigBuilder) BuildKongConfig() parser.KongConfigBuildingResult {
	return parser.KongConfigBuildingResult{
		KongState:           &kongstate.KongState{},
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

func TestKongClientUpdate_ConfigStatusIsAlwaysNotified(t *testing.T) {
	t.Parallel()

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
		kongClient             = setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder)
	)

	testCases := []struct {
		name                string
		gatewayFailure      bool
		konnectFailure      bool
		translationFailures bool
	}{
		{
			name:                "success",
			gatewayFailure:      false,
			konnectFailure:      false,
			translationFailures: false,
		},
		{
			name:                "gateway failure",
			gatewayFailure:      true,
			konnectFailure:      false,
			translationFailures: false,
		},
		{
			name:                "translation failures",
			gatewayFailure:      false,
			konnectFailure:      false,
			translationFailures: true,
		},
		{
			name:                "konnect failure",
			gatewayFailure:      false,
			konnectFailure:      true,
			translationFailures: false,
		},
		{
			name:                "both gateway and konnect failure",
			gatewayFailure:      true,
			konnectFailure:      true,
			translationFailures: false,
		},
		{
			name:                "translation failures and konnect failure",
			gatewayFailure:      false,
			konnectFailure:      true,
			translationFailures: true,
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
			require.True(t, statusQueue.WasNotified())
		})
	}
}

// setupTestKongClient creates a KongClient with mocked dependencies.
func setupTestKongClient(
	t *testing.T,
	updateStrategyResolver *mockUpdateStrategyResolver,
	clientsProvider mockGatewayClientsProvider,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
	configBuilder *mockKongConfigBuilder,
) *dataplane.KongClient {
	logger := logrus.New()
	timeout := time.Second
	ingressClass := "kong"
	diagnostic := util.ConfigDumpDiagnostic{}
	config := sendconfig.Config{}
	dbMode := "off"

	kongClient, err := dataplane.NewKongClient(
		logger,
		timeout,
		ingressClass,
		diagnostic,
		config,
		newFakeEventsRecorder(),
		dbMode,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
		configBuilder,
		store.NewCacheStores(),
	)
	require.NoError(t, err)
	return kongClient
}

func newFakeEventsRecorder() *record.FakeRecorder {
	eventRecorder := record.NewFakeRecorder(0)

	// Ingest events to unblock writing side.
	go func() {
		for range eventRecorder.Events {
		}
	}()

	return eventRecorder
}

func mustSampleGatewayClient(t *testing.T) *adminapi.Client {
	t.Helper()
	c, err := adminapi.NewTestClient(fmt.Sprintf("https://%s:8080", uuid.NewString()))
	require.NoError(t, err)
	return c
}

func mustSampleKonnectClient(t *testing.T) *adminapi.KonnectClient {
	t.Helper()

	c, err := gokong.NewClient(lo.ToPtr(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString())), &http.Client{})
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
