package dataplane

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/google/uuid"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configurationv1 "github.com/kong/kubernetes-configuration/v2/api/configuration/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/configfetcher"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckgen"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
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
	dbMode         dpconf.DBMode
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

// mockKongLastValidConfigFetcher is a mock implementation of FallbackConfigGenerator interface.
type mockFallbackConfigGenerator struct {
	GenerateResult store.CacheStores

	GenerateExcludingBrokenObjectsCalledWith   lo.Tuple2[store.CacheStores, []fallback.ObjectHash]
	GenerateBackfillingBrokenObjectsCalledWith lo.Tuple3[store.CacheStores, *store.CacheStores, []fallback.ObjectHash]
}

func newMockFallbackConfigGenerator() *mockFallbackConfigGenerator {
	return &mockFallbackConfigGenerator{}
}

func (m *mockFallbackConfigGenerator) GenerateExcludingBrokenObjects(
	stores store.CacheStores,
	hashes []fallback.ObjectHash,
) (store.CacheStores, fallback.GeneratedCacheMetadata, error) {
	m.GenerateExcludingBrokenObjectsCalledWith = lo.T2(stores, hashes)
	return m.GenerateResult, fallback.GeneratedCacheMetadata{}, nil
}

func (m *mockFallbackConfigGenerator) GenerateBackfillingBrokenObjects(
	currentStores store.CacheStores,
	lastValidStores *store.CacheStores,
	brokenObjects []fallback.ObjectHash,
) (store.CacheStores, fallback.GeneratedCacheMetadata, error) {
	m.GenerateBackfillingBrokenObjectsCalledWith = lo.T3(currentStores, lastValidStores, brokenObjects)
	return m.GenerateResult, fallback.GeneratedCacheMetadata{}, nil
}

func TestKongClientUpdate_AllExpectedClientsAreCalledAndErrorIsPropagated(t *testing.T) {
	var (
		ctx                = t.Context()
		testGatewayClients = []*adminapi.Client{
			mustSampleGatewayClient(t),
			mustSampleGatewayClient(t),
		}
	)

	testCases := []struct {
		name                 string
		gatewayClients       []*adminapi.Client
		withKonnectUpdater   bool
		errorOnUpdateForURLs []string
		expectError          bool
	}{
		{
			name:               "2 gateway clients and konnect with no errors",
			gatewayClients:     testGatewayClients,
			withKonnectUpdater: true,
			expectError:        false,
		},
		{
			name:                 "2 gateway clients with error on one of them",
			gatewayClients:       testGatewayClients,
			errorOnUpdateForURLs: []string{testGatewayClients[0].BaseRootURL()},
			expectError:          true,
		},
		{
			name:        "no clients at all",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientsProvider := &mockGatewayClientsProvider{
				gatewayClients: tc.gatewayClients,
			}
			updateStrategyResolver := mocks.NewUpdateStrategyResolver()
			for _, url := range tc.errorOnUpdateForURLs {
				updateStrategyResolver.ReturnErrorOnUpdate(url)
			}
			configChangeDetector := mocks.ConfigurationChangeDetector{
				// Always return true for ConfigurationChanged to trigger an update.
				ConfigurationChanged: true,
			}
			configBuilder := newMockKongConfigBuilder()
			kongRawStateGetter := &mockKongLastValidConfigFetcher{}
			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)

			// Set KonnectKongStateUpdater if requested.
			var konnectUpdater *mocks.KonnectKongStateUpdater
			if tc.withKonnectUpdater {
				konnectUpdater = &mocks.KonnectKongStateUpdater{}
				kongClient.SetKonnectKongStateUpdater(konnectUpdater)
			}

			err := kongClient.Update(ctx)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify that each gateway URL is called once.
			expectedURLsCalled := lo.SliceToMap(clientsProvider.GatewayClients(), func(c *adminapi.Client) (string, int) {
				return c.BaseRootURL(), 1
			})
			updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(t, expectedURLsCalled)
			// Verify that Konnect client is called eventually.
			if tc.withKonnectUpdater {
				require.Len(t, konnectUpdater.Calls(), 1, "expected Konnect updater to be called")
			}
		})
	}
}

func TestKongClientUpdate_WhenNoChangeInConfigNoClientGetsCalled(t *testing.T) {
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{
			mustSampleGatewayClient(t),
			mustSampleGatewayClient(t),
		},
	}
	updateStrategyResolver := mocks.NewUpdateStrategyResolver()

	// no change in config, we'll expect no update to be called
	configChangeDetector := mocks.ConfigurationChangeDetector{
		ConfigurationChanged: false,
	}
	configBuilder := newMockKongConfigBuilder()
	kongRawStateGetter := &mockKongLastValidConfigFetcher{}
	kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)

	ctx := t.Context()
	err := kongClient.Update(ctx)
	require.NoError(t, err)

	updateStrategyResolver.AssertNoUpdateCalled(t)
}

type mockConfigStatusQueue struct {
	gatewayConfigStatusNotifications []clients.GatewayConfigApplyStatus
	konnectConfigStatusNotifications []clients.KonnectConfigUploadStatus
	notifications                    []clients.ConfigStatus
	lock                             sync.RWMutex
}

func newMockConfigStatusQueue() *mockConfigStatusQueue {
	return &mockConfigStatusQueue{}
}

func (m *mockConfigStatusQueue) NotifyConfigStatus(_ context.Context, status clients.ConfigStatus) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.notifications = append(m.notifications, status)
}

func (m *mockConfigStatusQueue) NotifyGatewayConfigStatus(_ context.Context, status clients.GatewayConfigApplyStatus) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.gatewayConfigStatusNotifications = append(m.gatewayConfigStatusNotifications, status)
}

func (m *mockConfigStatusQueue) NotifyKonnectConfigStatus(_ context.Context, status clients.KonnectConfigUploadStatus) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.konnectConfigStatusNotifications = append(m.konnectConfigStatusNotifications, status)
}

func (m *mockConfigStatusQueue) GatewayConfigStatusNotifications() []clients.GatewayConfigApplyStatus {
	m.lock.RLock()
	defer m.lock.RUnlock()
	copied := make([]clients.GatewayConfigApplyStatus, len(m.gatewayConfigStatusNotifications))
	copy(copied, m.gatewayConfigStatusNotifications)
	return copied
}

func (m *mockConfigStatusQueue) KonnectConfigStatusNotifications() []clients.KonnectConfigUploadStatus {
	m.lock.RLock()
	defer m.lock.RUnlock()
	copied := make([]clients.KonnectConfigUploadStatus, len(m.konnectConfigStatusNotifications))
	copy(copied, m.konnectConfigStatusNotifications)
	return copied
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
	updateCacheCalls            []store.CacheStores

	// onlyFirstCallWithNoTranslationFailures is used to simulate a scenario where the first call to the
	// KongConfigBuilder has no translation failures, but subsequent calls do (e.g. to trigger translation failures in
	// fallback configuration).
	onlyFirstBuildCallWithNoTranslationFailures bool
	buildCalled                                 bool
}

func newMockKongConfigBuilder() *mockKongConfigBuilder {
	return &mockKongConfigBuilder{
		kongState: &kongstate.KongState{},
	}
}

func (p *mockKongConfigBuilder) BuildKongConfig() translator.KongConfigBuildingResult {
	if p.onlyFirstBuildCallWithNoTranslationFailures && !p.buildCalled {
		p.buildCalled = true
		return translator.KongConfigBuildingResult{
			KongState:           p.kongState,
			TranslationFailures: nil,
		}
	}
	return translator.KongConfigBuildingResult{
		KongState:           p.kongState,
		TranslationFailures: p.translationFailuresToReturn,
	}
}

func (p *mockKongConfigBuilder) UpdateCache(s store.CacheStores) {
	p.updateCacheCalls = append(p.updateCacheCalls, s)
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

func (p *mockKongConfigBuilder) CustomEntityTypes() []string {
	return nil
}

func (p *mockKongConfigBuilder) returnTranslationFailuresForAllButFirstCall(failures []failures.ResourceFailure) {
	p.onlyFirstBuildCallWithNoTranslationFailures = true
	p.translationFailuresToReturn = failures
}

func TestKongClientUpdate_ConfigStatusIsNotified(t *testing.T) {
	var (
		ctx               = t.Context()
		testGatewayClient = mustSampleGatewayClient(t)

		clientsProvider = &mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{testGatewayClient},
		}

		configChangeDetector = mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
		configBuilder        = newMockKongConfigBuilder()
	)

	testCases := []struct {
		name                 string
		gatewayFailuresCount int
		translationFailures  bool
		expectedNotification clients.GatewayConfigApplyStatus
	}{
		{
			name:                "success",
			translationFailures: false,
			expectedNotification: clients.GatewayConfigApplyStatus{
				TranslationFailuresOccurred: false,
				ApplyConfigFailed:           false,
			},
		},
		{
			name:                 "gateway failure",
			gatewayFailuresCount: 2,
			translationFailures:  false,
			expectedNotification: clients.GatewayConfigApplyStatus{
				TranslationFailuresOccurred: false,
				ApplyConfigFailed:           true,
			},
		},
		{
			name:                "translation failures",
			translationFailures: true,
			expectedNotification: clients.GatewayConfigApplyStatus{
				TranslationFailuresOccurred: true,
				ApplyConfigFailed:           false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				kongRawStateGetter     = &mockKongLastValidConfigFetcher{}
				updateStrategyResolver = mocks.NewUpdateStrategyResolver()
				statusQueue            = newMockConfigStatusQueue()
				kongClient             = setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)
			)

			kongClient.SetConfigStatusNotifier(statusQueue)
			for range tc.gatewayFailuresCount {
				updateStrategyResolver.ReturnErrorOnUpdate(testGatewayClient.BaseRootURL())
			}
			configBuilder.returnTranslationFailures(tc.translationFailures)

			_ = kongClient.Update(ctx)
			gatewayNotifications := statusQueue.GatewayConfigStatusNotifications()
			require.Len(t, gatewayNotifications, 1, "Should receive gateway configuration status right after update")

			notification := gatewayNotifications[0]
			assert.Equal(t, tc.expectedNotification, notification)
		})
	}
}

func TestKongClient_ApplyConfigurationEvents(t *testing.T) {
	testGatewayClient := mustSampleGatewayClient(t)
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{testGatewayClient},
	}
	updateStrategyResolver := mocks.NewUpdateStrategyResolver()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
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

			err := kongClient.Update(t.Context())
			require.NoError(t, err)

			if tc.expectEvents {
				require.NotEmpty(t, eventRecorder.Events())
			} else {
				require.Empty(t, eventRecorder.Events())
			}
		})
	}
}

func TestKongClient_KubernetesEvents(t *testing.T) {
	t.Setenv("POD_NAMESPACE", "test-namespace")
	t.Setenv("POD_NAME", "test-pod")

	ctx := t.Context()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	testIngress := helpers.WithTypeMeta(t, &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "obj-1",
			Namespace: "namespace",
		},
	})
	testService := helpers.WithTypeMeta(t, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "obj-2",
			Namespace: "namespace",
		},
	})

	testCases := []struct {
		name                                     string
		fallbackConfiguration                    bool
		translationFailures                      bool
		updateError                              bool
		entityErrors                             bool
		fallbackConfigurationUpdateError         bool
		fallbackConfigurationTranslationFailures bool
		expectError                              bool
		expectEmittingEvents                     []string
	}{
		{
			name:        "successful update",
			expectError: false,
			expectEmittingEvents: []string{
				"Normal KongConfigurationSucceeded",
			},
		},
		{
			name:                "translation failures",
			translationFailures: true,
			expectError:         false,
			expectEmittingEvents: []string{
				"Ingress: Warning KongConfigurationTranslationFailed",
				"Service: Warning KongConfigurationTranslationFailed",
				"Pod: Normal KongConfigurationSucceeded",
			},
		},
		{
			name:        "update error",
			updateError: true,
			expectError: true,
			expectEmittingEvents: []string{
				"Pod: Warning KongConfigurationApplyFailed",
			},
		},
		{
			name:         "update error with entity errors",
			updateError:  true,
			entityErrors: true,
			expectError:  true,
			expectEmittingEvents: []string{
				"Pod: Warning KongConfigurationApplyFailed",
				"Ingress: Warning KongConfigurationApplyFailed",
				"Service: Warning KongConfigurationApplyFailed",
			},
		},
		{
			name:                  "update error with entity errors, fallback configuration applied",
			fallbackConfiguration: true,
			updateError:           true,
			entityErrors:          true,
			expectError:           true,
			expectEmittingEvents: []string{
				"Pod: Warning KongConfigurationApplyFailed",
				"Ingress: Warning KongConfigurationApplyFailed",
				"Service: Warning KongConfigurationApplyFailed",
				"Pod: Normal FallbackKongConfigurationSucceeded",
			},
		},
		{
			name:                             "update error with entity errors, fallback configuration failures",
			fallbackConfiguration:            true,
			fallbackConfigurationUpdateError: true,
			updateError:                      true,
			entityErrors:                     true,
			expectError:                      true,
			expectEmittingEvents: []string{
				"Pod: Warning KongConfigurationApplyFailed",
				"Ingress: Warning KongConfigurationApplyFailed",
				"Service: Warning KongConfigurationApplyFailed",
				"Pod: Warning FallbackKongConfigurationApplyFailed",
				"Ingress: Warning FallbackKongConfigurationApplyFailed",
			},
		},
		{
			name:                                     "update error with entity errors, fallback translation failures",
			fallbackConfiguration:                    true,
			updateError:                              true,
			entityErrors:                             true,
			expectError:                              true,
			fallbackConfigurationTranslationFailures: true,
			expectEmittingEvents: []string{
				"Pod: Warning KongConfigurationApplyFailed",
				"Ingress: Warning KongConfigurationApplyFailed",
				"Service: Warning KongConfigurationApplyFailed",
				"Ingress: Warning FallbackKongConfigurationTranslationFailed",
				"Service: Warning FallbackKongConfigurationTranslationFailed",
				"Pod: Normal FallbackKongConfigurationSucceeded",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategyResolver := mocks.NewUpdateStrategyResolver()
			configBuilder := newMockKongConfigBuilder()
			eventRecorder := mocks.NewEventRecorder()
			lastValidConfigFetcher := &mockKongLastValidConfigFetcher{}
			testGatewayClient := mustSampleGatewayClient(t)
			clientsProvider := &mockGatewayClientsProvider{
				gatewayClients: []*adminapi.Client{testGatewayClient},
			}
			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, eventRecorder, lastValidConfigFetcher)
			kongClient.kongConfig.FallbackConfiguration = tc.fallbackConfiguration

			if tc.translationFailures {
				configBuilder.translationFailuresToReturn = []failures.ResourceFailure{
					lo.Must(failures.NewResourceFailure("some reason", testIngress)),
					lo.Must(failures.NewResourceFailure("some reason", testService)),
				}
			}
			if tc.updateError {
				if tc.entityErrors {
					updateStrategyResolver.ReturnSpecificErrorOnUpdate(testGatewayClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
						[]failures.ResourceFailure{
							lo.Must(failures.NewResourceFailure("violated constraint", testIngress)),
							lo.Must(failures.NewResourceFailure("violated constraint", testService)),
						},
						errors.New("error on update"),
					))
				} else {
					updateStrategyResolver.ReturnErrorOnUpdate(testGatewayClient.BaseRootURL())
				}
			}
			if tc.updateError && tc.fallbackConfigurationUpdateError {
				updateStrategyResolver.ReturnSpecificErrorOnUpdate(testGatewayClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
					[]failures.ResourceFailure{
						lo.Must(failures.NewResourceFailure("violated constraint", testIngress)),
					}, errors.New("error on update"),
				))
			}
			if tc.fallbackConfigurationTranslationFailures {
				configBuilder.returnTranslationFailuresForAllButFirstCall([]failures.ResourceFailure{
					lo.Must(failures.NewResourceFailure("some reason", testIngress)),
					lo.Must(failures.NewResourceFailure("some reason", testService)),
				})
			}

			err := kongClient.Update(ctx)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			emittedEvents := eventRecorder.Events()
			require.Len(t, emittedEvents, len(tc.expectEmittingEvents))
			for _, expectedEvent := range tc.expectEmittingEvents {
				containsExpectedEvent := lo.ContainsBy(emittedEvents, func(event string) bool {
					return strings.Contains(event, expectedEvent)
				})
				require.True(t, containsExpectedEvent, "expected event %q not found in %v", expectedEvent, eventRecorder.Events())
			}
		})
	}
}

func TestKongClient_EmptyConfigUpdate(t *testing.T) {
	var (
		ctx               = t.Context()
		testGatewayClient = mustSampleGatewayClient(t)

		clientsProvider = &mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{testGatewayClient},
		}

		updateStrategyResolver = mocks.NewUpdateStrategyResolver()
		configChangeDetector   = mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
		configBuilder          = newMockKongConfigBuilder()
		kongRawStateGetter     = &mockKongLastValidConfigFetcher{}
		kongClient             = setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)
	)

	t.Run("dbless", func(t *testing.T) {
		konnectKongStateUpdater := &mocks.KonnectKongStateUpdater{}
		kongClient.SetKonnectKongStateUpdater(konnectKongStateUpdater)

		kongClient.kongConfig.InMemory = true
		err := kongClient.Update(ctx)
		require.NoError(t, err)

		gwContent, ok := updateStrategyResolver.LastUpdatedContentForURL(testGatewayClient.BaseRootURL())
		require.True(t, ok)
		assert.Equal(t, gwContent.Content, &file.Content{
			FormatVersion: versions.DeckFileFormatVersion,
			Upstreams: []file.FUpstream{
				{
					Upstream: kong.Upstream{
						Name: lo.ToPtr(deckgen.StubUpstreamName),
					},
				},
			},
		}, "gateway content should have appended stub upstream")

		require.Len(t, konnectKongStateUpdater.Calls(), 1)
		konnectKongStateUpdaterCall := konnectKongStateUpdater.Calls()[0]
		assert.Equal(t, &kongstate.KongState{}, konnectKongStateUpdaterCall.KongState)
	})

	t.Run("db", func(t *testing.T) {
		konnectKongStateUpdater := &mocks.KonnectKongStateUpdater{}
		kongClient.SetKonnectKongStateUpdater(konnectKongStateUpdater)

		kongClient.kongConfig.InMemory = false
		err := kongClient.Update(ctx)
		require.NoError(t, err)

		gwContent, ok := updateStrategyResolver.LastUpdatedContentForURL(testGatewayClient.BaseRootURL())
		require.True(t, ok)
		require.True(t, deckgen.IsContentEmpty(gwContent.Content), "gateway content should be empty")

		require.Len(t, konnectKongStateUpdater.Calls(), 1)
		konnectKongStateUpdaterCall := konnectKongStateUpdater.Calls()[0]
		assert.Equal(t, &kongstate.KongState{}, konnectKongStateUpdaterCall.KongState)
	})
}

// setupTestKongClient creates a KongClient with mocked dependencies.
func setupTestKongClient(
	t *testing.T,
	updateStrategyResolver *mocks.UpdateStrategyResolver,
	clientsProvider *mockGatewayClientsProvider,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
	configBuilder *mockKongConfigBuilder,
	eventRecorder record.EventRecorder,
	kongRawStateGetter configfetcher.LastValidConfigFetcher,
) *KongClient {
	logger := zapr.NewLogger(zap.NewNop())
	timeout := time.Second
	config := sendconfig.Config{
		SanitizeKonnectConfigDumps: true,
	}

	if eventRecorder == nil {
		eventRecorder = mocks.NewEventRecorder()
	}

	cacheStores := store.NewCacheStores()
	kongClient, err := NewKongClient(
		logger,
		timeout,
		config,
		eventRecorder,
		dpconf.DBModeOff,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
		kongRawStateGetter,
		configBuilder,
		&cacheStores,
		newMockFallbackConfigGenerator(),
		mocks.MetricsRecorder{},
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

type mockKongLastValidConfigFetcher struct {
	kongRawState  *utils.KongRawState
	lastKongState *kongstate.KongState
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

func (cf *mockKongLastValidConfigFetcher) TryFetchingValidConfigFromGateways(context.Context, logr.Logger, []*adminapi.Client, []string) error {
	if cf.kongRawState != nil {
		cf.lastKongState = configfetcher.KongRawStateToKongState(cf.kongRawState)
	}
	return nil
}

func TestKongClientUpdate_FetchStoreAndPushLastValidConfig(t *testing.T) {
	var (
		ctx = t.Context()

		clientsProvider = &mockGatewayClientsProvider{
			gatewayClients: []*adminapi.Client{
				mustSampleGatewayClient(t),
				mustSampleGatewayClient(t),
			},
		}

		configChangeDetector = mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
		lastKongRawState     = &utils.KongRawState{
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
		updateErr     = sendconfig.NewUpdateErrorWithoutResponseBody(nil, errors.New("invalid request"))
	)
	configBuilder.kongState = newKongState

	testCases := []struct {
		name                   string
		translationFailures    bool
		errorOnGatewayFailures error
		gatewayFailuresCount   int
		lastValidKongRawState  *utils.KongRawState
		lastKongStatusHash     string
		expectedLastKongState  *kongstate.KongState
		errorsSize             int
	}{
		{
			name:                  "success, new fallback set",
			lastValidKongRawState: lastKongRawState,
			expectedLastKongState: newKongState,
			lastKongStatusHash:    "xyz",
		},
		{
			name:                   "no previous state, failure",
			errorOnGatewayFailures: updateErr,
			gatewayFailuresCount:   1,
			expectedLastKongState:  nil,
			errorsSize:             1,
			lastKongStatusHash:     sendconfig.WellKnownInitialHash,
		},
		{
			name:                   "previous state, failure, fallback pushed with success",
			errorOnGatewayFailures: updateErr,
			gatewayFailuresCount:   1,
			lastValidKongRawState:  lastKongRawState,
			expectedLastKongState:  lastKongState,
			errorsSize:             1,
			lastKongStatusHash:     "xyz",
		},
		{
			name:                   "previous state, failure, fallback pushed with failure",
			errorOnGatewayFailures: updateErr,
			gatewayFailuresCount:   2,
			lastValidKongRawState:  lastKongRawState,
			expectedLastKongState:  lastKongState,
			errorsSize:             3,
			lastKongStatusHash:     "xyz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategyResolver := mocks.NewUpdateStrategyResolver()
			for range tc.gatewayFailuresCount {
				updateStrategyResolver.ReturnSpecificErrorOnUpdate(clientsProvider.gatewayClients[0].BaseRootURL(), tc.errorOnGatewayFailures)
				updateStrategyResolver.ReturnSpecificErrorOnUpdate(clientsProvider.gatewayClients[1].BaseRootURL(), tc.errorOnGatewayFailures)
			}

			kongRawStateGetter := &mockKongLastValidConfigFetcher{
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

func TestKongClient_FallbackConfiguration_SuccessfulRecovery(t *testing.T) {
	ctx := t.Context()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	lastValidConfigFetcher := &mockKongLastValidConfigFetcher{}
	diagnosticsCh := make(chan diagnostics.ConfigDump, 10) // make it buffered to avoid blocking

	// We'll use KongConsumer as an example of a broken object, but it could be any supported type
	// for the purpose of this test as the fallback config generator is mocked anyway.
	validConsumer := someConsumer(t, "valid")
	brokenConsumer := someConsumer(t, "broken")
	originalCache := cacheStoresFromObjs(t, validConsumer, brokenConsumer)
	lastValidCache := cacheStoresFromObjs(t, validConsumer)

	testCases := []struct {
		name                                             string
		enableLastValidConfigFallback                    bool
		expectGenerateExcludingBrokenObjectsCalled       bool
		expectGenerateBackfillingBrokenObjectsCalledWith bool
	}{
		{
			name:                          "last valid config is disabled",
			enableLastValidConfigFallback: false,
			expectGenerateExcludingBrokenObjectsCalled: true,
		},
		{
			name:                          "last valid config is enabled",
			enableLastValidConfigFallback: true,
			expectGenerateBackfillingBrokenObjectsCalledWith: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateStrategyResolver := mocks.NewUpdateStrategyResolver()
			configBuilder := newMockKongConfigBuilder()
			fallbackConfigGenerator := newMockFallbackConfigGenerator()
			gwClient := mustSampleGatewayClient(t)
			clientsProvider := &mockGatewayClientsProvider{
				gatewayClients: []*adminapi.Client{gwClient},
			}
			kongClient, err := NewKongClient(
				zapr.NewLogger(zap.NewNop()),
				time.Second,
				sendconfig.Config{
					FallbackConfiguration:         true,
					UseLastValidConfigForFallback: tc.enableLastValidConfigFallback,
				},
				mocks.NewEventRecorder(),
				dpconf.DBModeOff,
				clientsProvider,
				updateStrategyResolver,
				configChangeDetector,
				lastValidConfigFetcher,
				configBuilder,
				&originalCache,
				fallbackConfigGenerator,
				mocks.MetricsRecorder{},
				WithDiagnosticsClient(diagnostics.Client{
					Configs: diagnosticsCh,
				}),
			)
			require.NoError(t, err)

			t.Log("Injecting the last valid cache snapshot to be used for recovery")
			kongClient.lastValidCacheSnapshot = &lastValidCache

			t.Log("Setting update strategy to return an error on the first call to trigger fallback configuration generation")
			updateStrategyResolver.ReturnSpecificErrorOnUpdate(gwClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
				[]failures.ResourceFailure{
					lo.Must(failures.NewResourceFailure("violated constraint", brokenConsumer)),
				},
				errors.New("error on update"),
			))

			t.Log("Setting the config builder to return KongState with the valid consumer only")
			configBuilder.kongState = &kongstate.KongState{
				Consumers: []kongstate.Consumer{
					{
						Consumer: kong.Consumer{
							Username: lo.ToPtr(validConsumer.Username),
						},
					},
				},
			}

			t.Log("Setting the fallback config generator to return a snapshot excluding the broken consumer")
			fallbackCacheStoresToBeReturned := cacheStoresFromObjs(t, validConsumer)
			fallbackConfigGenerator.GenerateResult = fallbackCacheStoresToBeReturned

			t.Log("Calling KongClient.Update")
			err = kongClient.Update(ctx)
			require.Error(t, err)

			t.Log("Verifying that the config builder cache was updated twice")
			require.Len(t, configBuilder.updateCacheCalls, 2,
				"expected cache to be updated with a snapshot twice: first with the initial cache snapshot, then with the fallback one")

			t.Log("Verifying that the first cache update contains both consumers")
			firstCacheUpdate := configBuilder.updateCacheCalls[0]
			require.NotEqual(t, originalCache, firstCacheUpdate, "expected cache to be updated with a new snapshot")
			_, hasConsumer, err := firstCacheUpdate.Consumer.Get(brokenConsumer)
			require.NoError(t, err)
			require.True(t, hasConsumer, "expected consumer to be in the first cache snapshot")

			if tc.expectGenerateExcludingBrokenObjectsCalled {
				t.Log("Verifying that the fallback config generator was called with the first cache snapshot and the broken object hash")
				expectedGenerateExcludingBrokenObjectsArgs := lo.T2(firstCacheUpdate, []fallback.ObjectHash{fallback.GetObjectHash(brokenConsumer)})
				require.Equal(t, expectedGenerateExcludingBrokenObjectsArgs, fallbackConfigGenerator.GenerateExcludingBrokenObjectsCalledWith,
					"expected fallback config generator to be called with the first cache snapshot and the broken object hash")

				require.Empty(t, fallbackConfigGenerator.GenerateBackfillingBrokenObjectsCalledWith)
			}
			if tc.expectGenerateBackfillingBrokenObjectsCalledWith {
				t.Log("Verifying that the fallback config generator was called with the first and last valid cache snapshots and the broken object hash")
				expectedGenerateBackfillingBrokenObjectsArgs := lo.T3(firstCacheUpdate, &lastValidCache, []fallback.ObjectHash{fallback.GetObjectHash(brokenConsumer)})
				require.Equal(t, expectedGenerateBackfillingBrokenObjectsArgs, fallbackConfigGenerator.GenerateBackfillingBrokenObjectsCalledWith,
					"expected fallback config generator to be called with the current and last valid cache snapshots and the broken object hash")
			}

			t.Log("Verifying that the second config builder cache update contains the fallback snapshot")
			secondCacheUpdate := configBuilder.updateCacheCalls[1]
			require.Equal(t, fallbackCacheStoresToBeReturned, secondCacheUpdate,
				"expected cache to be updated with the fallback snapshot on second call")

			t.Log("Verifying that the update strategy was called twice for gateway")
			updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
				t,
				map[string]int{
					gwClient.BaseRootURL(): 2,
				}, "expected update to be called twice: first with the initial config, then with the fallback one",
			)

			t.Log("Verifying that the last valid config is updated with the config excluding the broken consumer")
			lastValidConfig, ok := lastValidConfigFetcher.LastValidConfig()
			require.True(t, ok)
			require.Len(t, lastValidConfig.Consumers, 1)
			require.Equal(t, validConsumer.Username, *lastValidConfig.Consumers[0].Username)
		})
	}

	t.Log("Verifying that the last valid config is updated with the config excluding the broken consumer")
	lastValidConfig, _ := lastValidConfigFetcher.LastValidConfig()
	require.Len(t, lastValidConfig.Consumers, 1)
	require.Equal(t, validConsumer.Username, *lastValidConfig.Consumers[0].Username)

	t.Log("Verifying that the diagnostic server received a dump indicating that the broken consumer caused a problem")
	// The test will have pushed several successful configs that we don't care about into the diag buffer. This is a
	// silly hack to churn through those until we get to the successful fallback.
	var dump diagnostics.ConfigDump
	require.Eventually(t, func() bool {
		dump = <-diagnosticsCh
		return dump.Meta.Fallback
	}, time.Second, time.Nanosecond)

	// Once we have the fallback diagnostic dump, check to confirm that it was a successful fallback push.
	require.False(t, dump.Meta.Failed)
	require.True(t, dump.Meta.Fallback)
}

func TestKongClient_FallbackConfiguration_SkipsUpdateWhenInSync(t *testing.T) {
	ctx := t.Context()
	gwClient := mustSampleGatewayClient(t)
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{gwClient},
	}
	updateStrategyResolver := mocks.NewUpdateStrategyResolver()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	configBuilder := newMockKongConfigBuilder()
	lastValidConfigFetcher := &mockKongLastValidConfigFetcher{}
	fallbackConfigGenerator := newMockFallbackConfigGenerator()
	diagnosticsCh := make(chan diagnostics.ConfigDump, 10) // make it buffered to avoid blocking

	// We'll use KongConsumer as an example of an object, but it could be any supported type
	// for the purpose of this test as the fallback config generator is mocked anyway.
	originalCache := cacheStoresFromObjs(t, someConsumer(t, "valid"))
	kongClient, err := NewKongClient(
		zapr.NewLogger(zap.NewNop()),
		time.Second,
		sendconfig.Config{
			FallbackConfiguration: true,
		},
		mocks.NewEventRecorder(),
		dpconf.DBModeOff,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
		lastValidConfigFetcher,
		configBuilder,
		&originalCache,
		fallbackConfigGenerator,
		mocks.MetricsRecorder{},
		WithDiagnosticsClient(diagnostics.Client{
			Configs: diagnosticsCh,
		}),
	)
	require.NoError(t, err)

	t.Run("on first update clients are updated", func(t *testing.T) {
		t.Log("Calling KongClient.Update")
		require.NoError(t, kongClient.Update(ctx))

		t.Log("Verifying that the config builder cache was updated once")
		require.Len(t, configBuilder.updateCacheCalls, 1)

		t.Log("Verifying that the update strategy was called once for gateway")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(t, map[string]int{gwClient.BaseRootURL(): 1})
	})

	t.Run("without clients change, on second update clients are not updated", func(t *testing.T) {
		t.Log("Calling KongClient.Update again")
		require.NoError(t, kongClient.Update(ctx))

		t.Log("Verifying that the config builder cache was not updated when config was not changed")
		require.Len(t, configBuilder.updateCacheCalls, 1)

		t.Log("Verifying that the update strategy was not called again")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(t, map[string]int{gwClient.BaseRootURL(): 1})
	})

	newGwClient := mustSampleGatewayClient(t)
	t.Run("when new client is discovered, it is updated", func(t *testing.T) {
		t.Log("Injecting a new client to the provider")
		clientsProvider.gatewayClients = append(clientsProvider.gatewayClients, newGwClient)

		t.Log("Calling KongClient.Update again")
		require.NoError(t, kongClient.Update(ctx))

		t.Log("Verifying that the config builder cache is not updated as there was no config change")
		require.Len(t, configBuilder.updateCacheCalls, 1)

		t.Log("Verifying that the update strategies were called for the client that was added")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
			t,
			map[string]int{
				gwClient.BaseRootURL():    2, // First series of updates + Second series of updates
				newGwClient.BaseRootURL(): 1, // Second series of updates only
			})
	})

	t.Run("when generating fallback, all clients are updated", func(t *testing.T) {
		t.Log("Changing configuration")
		require.NoError(t, originalCache.Add(someConsumer(t, "broken"))) // Add a consumer to cache to change the cache hash.

		t.Log("Setting update strategy to return an error on the first call to trigger fallback configuration generation")
		updateStrategyResolver.ReturnSpecificErrorOnUpdate(gwClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
			[]failures.ResourceFailure{
				lo.Must(failures.NewResourceFailure("violated constraint", someConsumer(t, "invalid"))),
			},
			errors.New("error on update"),
		))

		t.Log("Calling KongClient.Update")
		require.Error(t, kongClient.Update(ctx))

		t.Log("Verifying that the update strategy was called again for all gateways")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
			t,
			map[string]int{
				gwClient.BaseRootURL():    4, // First series of updates + Second series of updates + rejected/fallback
				newGwClient.BaseRootURL(): 3, // Second series of updates + rejected/fallback
			})
	})

	anotherNewGwClient := mustSampleGatewayClient(t)
	t.Run("when fallback was used before and config is still broken, after discovering a new client, all clients are updated", func(t *testing.T) {
		t.Log("Adding a new client to the provider")
		clientsProvider.gatewayClients = append(clientsProvider.gatewayClients, anotherNewGwClient)
		updateStrategyResolver.ReturnSpecificErrorOnUpdate(gwClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
			[]failures.ResourceFailure{
				lo.Must(failures.NewResourceFailure("violated constraint", someConsumer(t, "invalid"))),
			},
			errors.New("error on update"),
		))

		t.Log("Calling KongClient.Update again")
		require.Error(t, kongClient.Update(ctx))

		t.Log("Verifying that the update strategy was called again for all gateways")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
			t,
			map[string]int{
				gwClient.BaseRootURL():    6, // First series of updates + Second series of updates + rejected/fallback * 2
				newGwClient.BaseRootURL(): 5, // Second series of updates + rejected/fallback * 2
			})
	})

	t.Run("after fallback, when new config is correct, all clients are updated", func(t *testing.T) {
		t.Log("Changing configuration")
		require.NoError(t, originalCache.Consumer.Add(someConsumer(t, "valid"))) // Add a consumer to cache to change the cache hash.

		t.Log("Calling KongClient.Update")
		require.NoError(t, kongClient.Update(ctx))

		t.Log("Verifying that the update strategy was called again for all gateways")
		updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
			t,
			map[string]int{
				gwClient.BaseRootURL():    7, // First series of updates + Second series of updates + rejected/fallback * 2 + third update
				newGwClient.BaseRootURL(): 6, // Second series of updates + rejected/fallback + rejected/fallback * 2 + third update
			})
	})
}

func TestKongClient_FallbackConfiguration_FailedRecovery(t *testing.T) {
	ctx := t.Context()
	gwClient := mustSampleGatewayClient(t)
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{gwClient},
	}
	updateStrategyResolver := mocks.NewUpdateStrategyResolver()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	configBuilder := newMockKongConfigBuilder()
	lastValidConfigFetcher := &mockKongLastValidConfigFetcher{}
	fallbackConfigGenerator := newMockFallbackConfigGenerator()
	diagnosticsCh := make(chan diagnostics.ConfigDump, 10) // make it buffered to avoid blocking

	// We'll use KongConsumer as an example of a broken object, but it could be any supported type
	// for the purpose of this test as the fallback config generator is mocked anyway.
	brokenConsumer := someConsumer(t, "broken")
	originalCache := cacheStoresFromObjs(t, brokenConsumer)
	kongClient, err := NewKongClient(
		zapr.NewLogger(zap.NewNop()),
		time.Second,
		sendconfig.Config{
			FallbackConfiguration: true,
		},
		mocks.NewEventRecorder(),
		dpconf.DBModeOff,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
		lastValidConfigFetcher,
		configBuilder,
		&originalCache,
		fallbackConfigGenerator,
		mocks.MetricsRecorder{},
		WithDiagnosticsClient(diagnostics.Client{
			Configs: diagnosticsCh,
		}),
	)
	require.NoError(t, err)

	t.Log("Setting update strategy to return an error on the first call to trigger fallback configuration generation")
	updateStrategyResolver.ReturnSpecificErrorOnUpdate(gwClient.BaseRootURL(), sendconfig.NewUpdateErrorWithoutResponseBody(
		[]failures.ResourceFailure{
			lo.Must(failures.NewResourceFailure("violated constraint", brokenConsumer)),
		},
		errors.New("error on update"),
	))

	t.Log("Setting update strategy to return an error on the second call (fallback) to trigger a failed recovery")
	updateStrategyResolver.ReturnErrorOnUpdate(gwClient.BaseRootURL())

	t.Log("Calling KongClient.Update")
	err = kongClient.Update(ctx)
	require.Error(t, err)

	t.Log("Verifying that the update strategy was called twice for gateway")
	updateStrategyResolver.AssertUpdateCalledForURLsWithGivenCount(
		t,
		map[string]int{gwClient.BaseRootURL(): 2},
		"expected update to be called twice: first with the initial config, then with the fallback one")

	t.Log("Verifying that the last valid config is empty")
	_, hasLastValidConfig := lastValidConfigFetcher.LastValidConfig()
	require.False(t, hasLastValidConfig, "expected no last valid config to be stored as no successful recovery happened")

	t.Log("Verifying that the diagnostic server received a dump indicating that the broken consumer caused a problem")
	// The test will have pushed several successful configs that we don't care about into the diag buffer. This is a
	// silly hack to churn through those until we get to the failed fallback.
	var dump diagnostics.ConfigDump
	require.Eventually(t, func() bool {
		dump = <-diagnosticsCh
		return dump.Meta.Fallback
	}, time.Second, time.Nanosecond)

	// Once we have the fallback diagnostic dump, check to confirm that it was a failed fallback push.
	require.True(t, dump.Meta.Failed)
}

func TestKongClient_LastValidCacheSnapshot(t *testing.T) {
	var (
		ctx                     = t.Context()
		updateStrategyResolver  = mocks.NewUpdateStrategyResolver()
		configChangeDetector    = mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
		configBuilder           = newMockKongConfigBuilder()
		lastValidConfigFetcher  = &mockKongLastValidConfigFetcher{}
		originalCache           = cacheStoresFromObjs(t)
		fallbackConfigGenerator = newMockFallbackConfigGenerator()
	)

	testCases := []struct {
		name                                 string
		fallbackConfigurationFeatureEnabled  bool
		useLastValidConfigForFallbackEnabled bool
		expectLastValidCacheSnapshotToBeSet  bool
	}{
		{
			name:                                 "FallbackConfiguration=true, UseLastValidConfigForFallback=false",
			fallbackConfigurationFeatureEnabled:  true,
			useLastValidConfigForFallbackEnabled: false,
			expectLastValidCacheSnapshotToBeSet:  false,
		},
		{
			name:                                 "FallbackConfiguration=true, UseLastValidConfigForFallback=true",
			fallbackConfigurationFeatureEnabled:  true,
			useLastValidConfigForFallbackEnabled: true,
			expectLastValidCacheSnapshotToBeSet:  true,
		},
		{
			name:                                 "FallbackConfiguration=false, UseLastValidConfigForFallback=false",
			fallbackConfigurationFeatureEnabled:  false,
			useLastValidConfigForFallbackEnabled: false,
			expectLastValidCacheSnapshotToBeSet:  false,
		},
		{
			name:                                 "FallbackConfiguration=false, UseLastValidConfigForFallback=true",
			fallbackConfigurationFeatureEnabled:  false,
			useLastValidConfigForFallbackEnabled: true,
			expectLastValidCacheSnapshotToBeSet:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGatewayClient := mustSampleGatewayClient(t)
			clientsProvider := &mockGatewayClientsProvider{
				gatewayClients: []*adminapi.Client{testGatewayClient},
			}

			kongClient, err := NewKongClient(
				zapr.NewLogger(zap.NewNop()),
				time.Second,
				sendconfig.Config{
					FallbackConfiguration:         tc.fallbackConfigurationFeatureEnabled,
					UseLastValidConfigForFallback: tc.useLastValidConfigForFallbackEnabled,
				},
				mocks.NewEventRecorder(),
				dpconf.DBModeOff,
				clientsProvider,
				updateStrategyResolver,
				configChangeDetector,
				lastValidConfigFetcher,
				configBuilder,
				&originalCache,
				fallbackConfigGenerator,
				mocks.MetricsRecorder{},
			)
			require.NoError(t, err)

			require.Empty(t, kongClient.lastValidCacheSnapshot, "expected last valid cache snapshot to be empty")

			err = kongClient.Update(ctx)
			require.NoError(t, err)

			lastValid := kongClient.lastValidCacheSnapshot
			if tc.expectLastValidCacheSnapshotToBeSet {
				require.NotEmpty(t, lastValid, "expected last valid cache snapshot to be set after successful update")
			} else {
				require.Empty(t, lastValid, "expected last valid cache snapshot to remain empty")
			}
		})
	}
}

func cacheStoresFromObjs(t *testing.T, objs ...runtime.Object) store.CacheStores {
	for i := range objs {
		obj := objs[i].(client.Object)
		obj = helpers.WithTypeMeta(t, obj)
		objs[i] = obj
	}
	s, err := store.NewCacheStoresFromObjs(objs...)
	require.NoError(t, err)
	return s
}

func TestKongClient_ConfigDumpSanitization(t *testing.T) {
	clientsProvider := &mockGatewayClientsProvider{
		gatewayClients: []*adminapi.Client{
			mustSampleGatewayClient(t),
		},
	}
	updateStrategyResolver := mocks.NewUpdateStrategyResolver()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	configBuilder := newMockKongConfigBuilder()
	kongRawStateGetter := &mockKongLastValidConfigFetcher{}

	const testPrivateKey = "private-key-string"
	configBuilder.kongState = &kongstate.KongState{
		Certificates: []kongstate.Certificate{
			{
				Certificate: kong.Certificate{
					ID:  kong.String("new_cert"),
					Key: kong.String(testPrivateKey), // This should be redacted.
				},
			},
		},
	}
	kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector, configBuilder, nil, kongRawStateGetter)

	testCases := []struct {
		name                  string
		dumpsIncludeSensitive bool
		expectSanitizedDump   bool
	}{
		{
			name:                  "when DumpsIncludeSensitive is true, expect no sanitization",
			dumpsIncludeSensitive: true,
			expectSanitizedDump:   false,
		},
		{
			name:                  "when DumpsIncludeSensitive is false, expect sanitization",
			dumpsIncludeSensitive: false,
			expectSanitizedDump:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diagnosticsCh := make(chan diagnostics.ConfigDump, 1) // make it buffered to avoid blocking
			kongClient.diagnostic = diagnostics.Client{
				Configs:               diagnosticsCh,
				DumpsIncludeSensitive: tc.dumpsIncludeSensitive,
			}
			ctx := t.Context()
			err := kongClient.Update(ctx)
			require.NoError(t, err)

			dump := <-diagnosticsCh
			require.NotNil(t, dump.Config)
			require.Len(t, dump.Config.Certificates, 1)
			dumpedCert := dump.Config.Certificates[0]
			if tc.expectSanitizedDump {
				require.Equal(t, "{vault://redacted-value}", *dumpedCert.Key)
			} else {
				require.Equal(t, testPrivateKey, *dumpedCert.Key)
			}
		})
	}
}

func TestKongClient_RecoveringFromGatewaySyncError(t *testing.T) {
	ctx := t.Context()
	configChangeDetector := mocks.ConfigurationChangeDetector{ConfigurationChanged: true}
	fallbackConfigGenerator := newMockFallbackConfigGenerator()
	originalCache := cacheStoresFromObjs(t)

	testCases := []struct {
		name                                     string
		errorsFromGateways                       []error
		hasLastValidConfig                       bool
		expectRecoveryByGeneratingFallbackConfig bool
		expectRecoveryByApplyingLastValidConfig  bool
	}{
		{
			name: "one of gateways returns UpdateError with entities",
			errorsFromGateways: []error{
				sendconfig.NewUpdateErrorWithoutResponseBody(
					[]failures.ResourceFailure{
						lo.Must(failures.NewResourceFailure("violated constraint", someConsumer(t, "broken"))),
					},
					errors.New("error on update"),
				),
				nil,
			},
			expectRecoveryByGeneratingFallbackConfig: true,
		},
		{
			name: "one of gateways returns UpdateError without entities, has last valid config",
			errorsFromGateways: []error{
				sendconfig.NewUpdateErrorWithoutResponseBody(nil, errors.New("error on update")),
				nil,
			},
			hasLastValidConfig:                      true,
			expectRecoveryByApplyingLastValidConfig: true,
		},
		{
			name: "one of gateways returns UpdateError without entities, no last valid config",
			errorsFromGateways: []error{
				sendconfig.NewUpdateErrorWithoutResponseBody(nil, errors.New("error on update")),
				nil,
			},
			hasLastValidConfig:                       false,
			expectRecoveryByGeneratingFallbackConfig: false,
			expectRecoveryByApplyingLastValidConfig:  false,
		},
		{
			name: "one of gateways returns unexpected error",
			errorsFromGateways: []error{
				errors.New("unexpected error on update"),
				nil,
			},
			hasLastValidConfig:                       true,
			expectRecoveryByGeneratingFallbackConfig: false,
			expectRecoveryByApplyingLastValidConfig:  false,
		},
		{
			name: "one gateway returns UpdateError, another one an unexpected error",
			errorsFromGateways: []error{
				sendconfig.NewUpdateErrorWithoutResponseBody(
					[]failures.ResourceFailure{
						lo.Must(failures.NewResourceFailure("violated constraint", someConsumer(t, "broken"))),
					},
					errors.New("error on update"),
				),
				errors.New("unexpected error on update"),
				nil,
			},
			hasLastValidConfig:                       true,
			expectRecoveryByGeneratingFallbackConfig: true,
			expectRecoveryByApplyingLastValidConfig:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Preparing %d gateway clients", len(tc.errorsFromGateways))
			updateStrategyResolver := mocks.NewUpdateStrategyResolver()
			gwClients := make([]*adminapi.Client, len(tc.errorsFromGateways))
			for i := range gwClients {
				gwClients[i] = mustSampleGatewayClient(t)
				updateStrategyResolver.ReturnSpecificErrorOnUpdate(gwClients[i].BaseRootURL(), tc.errorsFromGateways[i])
			}
			clientsProvider := &mockGatewayClientsProvider{
				gatewayClients: gwClients,
			}

			lastValidConfigFetcher := &mockKongLastValidConfigFetcher{}
			if tc.hasLastValidConfig {
				t.Logf("Setting last valid config to contain a consumer with username 'last-valid'")
				lastValidConfigFetcher.lastKongState = &kongstate.KongState{
					Consumers: []kongstate.Consumer{
						{
							Consumer: kong.Consumer{
								Username: lo.ToPtr("last-valid"),
							},
						},
					},
				}
			}

			t.Logf("Preparing config builder with a consumer with username 'fallback'")
			configBuilder := newMockKongConfigBuilder()
			configBuilder.kongState = &kongstate.KongState{
				Consumers: []kongstate.Consumer{
					{
						Consumer: kong.Consumer{
							Username: lo.ToPtr("fallback"),
						},
					},
				},
			}

			kongClient, err := NewKongClient(
				zapr.NewLogger(zap.NewNop()),
				time.Second,
				sendconfig.Config{
					FallbackConfiguration: true,
				},
				mocks.NewEventRecorder(),
				dpconf.DBModeOff,
				clientsProvider,
				updateStrategyResolver,
				configChangeDetector,
				lastValidConfigFetcher,
				configBuilder,
				&originalCache,
				fallbackConfigGenerator,
				mocks.MetricsRecorder{},
			)
			require.NoError(t, err)

			err = kongClient.Update(ctx)
			require.Error(t, err)

			expectedUpdatedURLs := lo.Map(gwClients, func(c *adminapi.Client, _ int) string {
				return c.BaseRootURL()
			})
			if tc.expectRecoveryByGeneratingFallbackConfig || tc.expectRecoveryByApplyingLastValidConfig {
				// In case of any recovery method, we expect the update to be called twice for each gateway.
				expectedUpdatedURLs = slices.Concat(expectedUpdatedURLs, expectedUpdatedURLs)
			}
			t.Logf("Ensuring that the update strategy was called %d times", len(expectedUpdatedURLs))
			updateStrategyResolver.AssertUpdateCalledForURLs(t, expectedUpdatedURLs)

			expectedContent := func(consumerUsername string) *file.Content {
				return &file.Content{
					FormatVersion: "3.0",
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: lo.ToPtr(consumerUsername),
							},
						},
					},
				}
			}
			receivedContent, ok := updateStrategyResolver.LastUpdatedContentForURL(expectedUpdatedURLs[0])
			require.True(t, ok)
			if tc.expectRecoveryByApplyingLastValidConfig {
				t.Log("Verifying that the last valid config was applied")
				require.Equal(t, expectedContent("last-valid"), receivedContent.Content)
			}
			if tc.expectRecoveryByGeneratingFallbackConfig {
				t.Log("Verifying that the fallback config was generated and applied")
				require.Equal(t, expectedContent("fallback"), receivedContent.Content)
			}
		})
	}
}

func someConsumer(t *testing.T, name string) *configurationv1.KongConsumer {
	return helpers.WithTypeMeta(t, &configurationv1.KongConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "name",
			Namespace: "namespace",
			Annotations: map[string]string{
				annotations.IngressClassKey: annotations.DefaultIngressClass,
			},
			UID: k8stypes.UID(uuid.NewString()),
		},
		Username: name,
	})
}
