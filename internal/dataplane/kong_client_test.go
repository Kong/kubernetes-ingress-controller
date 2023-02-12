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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlfake "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
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

type konnectClient struct {
	isKonnect bool
}

func (c konnectClient) IsKonnect() bool {
	return c.isKonnect
}

func TestHandleSendToClientResult(t *testing.T) {
	const testSHA = "2110454484b88378619111aab0d8a8b8d0ecad5c0ad1120a19810c965f8652dd"
	testError := errors.New("sending to client failure")

	testCases := []struct {
		name      string
		isKonnect bool
		inputSHA  string
		inputErr  error

		expectedErr error
		expectedSHA string
	}{
		{
			name:     "no error, sha is passed",
			inputSHA: testSHA,

			expectedSHA: testSHA,
		},
		{
			name:     "error is passed",
			inputSHA: testSHA,
			inputErr: testError,

			expectedErr: testError,
		},
		{
			name:      "konnect - error is ignored",
			isKonnect: true,
			inputSHA:  testSHA,
			inputErr:  testError,

			expectedErr: nil,
			expectedSHA: "",
		},
		{
			name:      "konnect - no error, sha is passed",
			isKonnect: true,
			inputSHA:  testSHA,

			expectedSHA: testSHA,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := konnectClient{isKonnect: tc.isKonnect}
			resultSHA, err := dataplane.HandleSendToClientResult(c, logrus.New(), tc.inputSHA, tc.inputErr)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedSHA, resultSHA)
		})
	}
}

// mockGatewayClientsProvider is a mock implementation of dataplane.AdminAPIClientsProvider.
type mockGatewayClientsProvider struct {
	gatewayClients []*adminapi.Client
	konnectClient  *adminapi.Client
}

func (f mockGatewayClientsProvider) AllClients() []*adminapi.Client {
	all := make([]*adminapi.Client, len(f.gatewayClients))
	copy(all, f.gatewayClients)
	if f.konnectClient != nil {
		all = append(all, f.konnectClient)
	}
	return all
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
func (f *mockUpdateStrategyResolver) returnErrorOnUpdate(url string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.shouldReturnErrorOnUpdate[url] = struct{}{}
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

func (m *mockUpdateStrategy) Update(context.Context, *file.Content) (
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

// mockConfigurationChangeDetector is a mock implementation of sendconfig.ConfigurationChangeDetector.
type mockConfigurationChangeDetector struct {
	hasConfigurationChanged bool
}

func (m mockConfigurationChangeDetector) HasConfigurationChanged(
	context.Context, []byte, []byte, *file.Content, sendconfig.KonnectAwareClient, sendconfig.StatusClient,
) (bool, error) {
	return m.hasConfigurationChanged, nil
}

func TestKongClientUpdate_AllExpectedClientsAreCalled(t *testing.T) {
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
		konnectClient        *adminapi.Client
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
				updateStrategyResolver.returnErrorOnUpdate(url)
			}
			// always return true for HasConfigurationChanged to trigger an update
			configChangeDetector := mockConfigurationChangeDetector{hasConfigurationChanged: true}

			kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector)

			err := kongClient.Update(ctx)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			allExpectedURLs := mapClientsToUrls(clientsProvider.AllClients())
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

	kongClient := setupTestKongClient(t, updateStrategyResolver, clientsProvider, configChangeDetector)

	ctx := context.Background()
	err := kongClient.Update(ctx)
	require.NoError(t, err)

	updateStrategyResolver.assertNoUpdateCalled()
}

// setupTestKongClient creates a KongClient with mocked dependencies.
func setupTestKongClient(
	t *testing.T,
	updateStrategyResolver *mockUpdateStrategyResolver,
	clientsProvider mockGatewayClientsProvider,
	configChangeDetector sendconfig.ConfigurationChangeDetector,
) *dataplane.KongClient {
	logger := logrus.New()
	timeout := time.Second
	ingressClass := "kong"
	diagnostic := util.ConfigDumpDiagnostic{}
	config := sendconfig.Config{}
	eventRecorder := record.NewFakeRecorder(0)
	dbMode := "off"
	client := ctrlfake.NewClientBuilder().Build()

	kongClient, err := dataplane.NewKongClient(
		logger,
		timeout,
		ingressClass,
		diagnostic,
		config,
		eventRecorder,
		dbMode,
		client,
		clientsProvider,
		updateStrategyResolver,
		configChangeDetector,
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

func mustSampleKonnectClient(t *testing.T) *adminapi.Client {
	t.Helper()

	c, err := gokong.NewClient(lo.ToPtr(fmt.Sprintf("https://%s.konghq.tech", uuid.NewString())), &http.Client{})
	require.NoError(t, err)

	rgID := uuid.NewString()
	return adminapi.NewKonnectClient(c, rgID)
}

func mapClientsToUrls(clients []*adminapi.Client) []string {
	return lo.Map(clients, func(c *adminapi.Client, _ int) string {
		return c.BaseRootURL()
	})
}
