package clients_test

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/clients"
	"github.com/kong/kubernetes-ingress-controller/v3/test/mocks"
)

type readinessCheckCall struct {
	AlreadyCreatedURLs []string
	PendingURLs        []string
}

type mockReadinessChecker struct {
	nextResult clients.ReadinessCheckResult
	lastCall   mo.Option[readinessCheckCall]
	callsCount int
	lock       sync.RWMutex
}

func (m *mockReadinessChecker) CheckReadiness(
	_ context.Context,
	alreadyCreatedClients []clients.AlreadyCreatedClient,
	pendingClients []adminapi.DiscoveredAdminAPI,
) clients.ReadinessCheckResult {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.callsCount++
	m.lastCall = mo.Some(readinessCheckCall{
		AlreadyCreatedURLs: lo.Map(alreadyCreatedClients, func(c clients.AlreadyCreatedClient, _ int) string {
			return c.BaseRootURL()
		}),
		PendingURLs: lo.Map(pendingClients, func(c adminapi.DiscoveredAdminAPI, _ int) string {
			return c.Address
		}),
	})
	return m.nextResult
}

func (m *mockReadinessChecker) LetChecksReturn(result clients.ReadinessCheckResult) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.nextResult = result
}

func (m *mockReadinessChecker) LastCall() (readinessCheckCall, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if call, ok := m.lastCall.Get(); ok {
		return call, true
	}
	return readinessCheckCall{}, false
}

func (m *mockReadinessChecker) CallsCount() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.callsCount
}

func intoTurnedReady(urls ...string) []*adminapi.Client {
	return lo.Map(urls, func(url string, _ int) *adminapi.Client {
		return lo.Must(adminapi.NewTestClient(url))
	})
}

func intoTurnedPending(urls ...string) []adminapi.DiscoveredAdminAPI {
	return lo.Map(urls, func(url string, _ int) adminapi.DiscoveredAdminAPI {
		return testDiscoveredAdminAPI(url)
	})
}

func TestAdminAPIClientsManager_OnNotifyClientsAreUpdatedAccordingly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zapr.NewLogger(zap.NewNop())
	readinessChecker := &mockReadinessChecker{}
	initialClient, err := adminapi.NewTestClient("https://localhost:8083")
	require.NoError(t, err)
	manager, err := clients.NewAdminAPIClientsManager(
		ctx,
		logger,
		[]*adminapi.Client{initialClient},
		readinessChecker,
	)
	require.NoError(t, err)
	require.NotNil(t, manager)
	manager.Run()
	<-manager.Running()

	requireClientsMatchEventually := func(t *testing.T, c *clients.AdminAPIClientsManager, addresses []string, args ...any) {
		require.Eventually(t, func() bool {
			clientAddresses := lo.Map(c.GatewayClients(), func(cl *adminapi.Client, _ int) string {
				return cl.BaseRootURL()
			})
			return slices.Equal(addresses, clientAddresses)
		}, time.Second, time.Millisecond, args...)
	}

	requireClientsMatchEventually(t, manager, []string{initialClient.BaseRootURL()},
		"initially there should be the initial client")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL1)})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})
	requireClientsMatchEventually(t, manager, []string{testURL1},
		"after notifying about a new address we should get 1 client eventually")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})
	requireClientsMatchEventually(t, manager, []string{testURL1},
		"after notifying the same address there's no update in clients")

	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1), testDiscoveredAdminAPI(testURL2)})
	requireClientsMatchEventually(t, manager, []string{testURL1},
		"after notifying new address set including the old already existing one but new one not yet ready we get just the old one")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL2)})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1), testDiscoveredAdminAPI(testURL2)})
	requireClientsMatchEventually(t, manager, []string{testURL1, testURL2},
		"after notifying new address set including the old already existing one and new one turning ready we get both the old and the new")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1), testDiscoveredAdminAPI(testURL2)})
	requireClientsMatchEventually(t, manager, []string{testURL1, testURL2},
		"after notifying again with the same set of URLs should not change the existing URLs")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedPending: intoTurnedPending(testURL2)})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1), testDiscoveredAdminAPI(testURL2)})
	requireClientsMatchEventually(t, manager, []string{testURL1},
		"after notifying the same address set with one turning pending, we get only one client")

	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{})
	manager.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})
	requireClientsMatchEventually(t, manager, []string{testURL1},
		"notifying again with just one URL should decrease the set of URLs to just this one")

	manager.Notify([]adminapi.DiscoveredAdminAPI{})
	requireClientsMatchEventually(t, manager, []string{})

	cancel()
	require.NotPanics(t, func() { manager.Notify([]adminapi.DiscoveredAdminAPI{}) }, "notifying about new clients after manager has been shut down shouldn't panic")
}

func TestNewAdminAPIClientsManager_NoInitialClientsDisallowed(t *testing.T) {
	_, err := clients.NewAdminAPIClientsManager(
		context.Background(),
		zapr.NewLogger(zap.NewNop()),
		nil,
		&mockReadinessChecker{},
	)
	require.ErrorContains(t, err, "at least one initial client must be provided")
}

func TestAdminAPIClientsManager_NotRunningNotifyLoop(t *testing.T) {
	t.Parallel()

	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	m, err := clients.NewAdminAPIClientsManager(
		context.Background(),
		zapr.NewLogger(zap.NewNop()),
		[]*adminapi.Client{testClient},
		&mockReadinessChecker{},
	)
	require.NoError(t, err)

	select {
	case <-m.Running():
		t.Error("expected manager to not run without explicitly running it with Run method")
	case <-time.After(time.Millisecond * 100):
	}
}

func TestAdminAPIClientsManager_Clients(t *testing.T) {
	t.Parallel()

	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	m, err := clients.NewAdminAPIClientsManager(
		context.Background(),
		zapr.NewLogger(zap.NewNop()),
		[]*adminapi.Client{testClient},
		&mockReadinessChecker{},
	)
	require.NoError(t, err)
	require.Len(t, m.GatewayClients(), 1, "expecting one initial client")
	require.Equal(t, m.GatewayClientsCount(), 1, "expecting one initial client")
	require.Len(t, m.GatewayClientsToConfigure(), 1, "Expecting one initial client")

	konnectTestClient := &adminapi.KonnectClient{}
	m.SetKonnectClient(konnectTestClient)
	require.Len(t, m.GatewayClients(), 1, "konnect client should not be returned from GatewayClients")
	require.Equal(t, m.GatewayClientsCount(), 1, "konnect client should not be counted in GatewayClientsCount")
	require.Equal(t, konnectTestClient, m.KonnectClient(), "konnect client should be returned from KonnectClient")
}

func TestAdminAPIClientsManager_Clients_DBMode(t *testing.T) {
	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	testClient2, err := adminapi.NewTestClient("localhost:8081")
	require.NoError(t, err)
	initialClients := []*adminapi.Client{testClient, testClient2}
	require.NoError(t, err)

	m, err := clients.NewAdminAPIClientsManager(
		context.Background(),
		zapr.NewLogger(zap.NewNop()),
		initialClients,
		&mockReadinessChecker{},
	)
	require.NoError(t, err)
	m = m.WithDBMode("postgres")

	clients := m.GatewayClients()
	require.Len(t, clients, 2, "Expecting 2 clients returned with DB mode")

	configureClients := m.GatewayClientsToConfigure()
	require.Len(t, configureClients, 1, "Expecting 1 client to configure")
	require.Truef(t, lo.ContainsBy(initialClients, func(c *adminapi.Client) bool {
		return c.BaseRootURL() == configureClients[0].BaseRootURL()
	}), "Client's address %s should be in initial clients")

	require.Equal(t, m.GatewayClientsCount(), 2, "Expecting 2 initial clients")
}

func TestAdminAPIClientsManager_SubscribeToGatewayClientsChanges(t *testing.T) {
	t.Parallel()

	readinessChecker := &mockReadinessChecker{}
	testClient, err := adminapi.NewTestClient("http://localhost:8080")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	m, err := clients.NewAdminAPIClientsManager(
		ctx,
		zapr.NewLogger(zap.NewNop()),
		[]*adminapi.Client{testClient},
		readinessChecker)

	require.NoError(t, err)

	t.Run("no notify loop running should return false when subscribing", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.Nil(t, ch)
		require.Falsef(t, ok, "expected no subscription to be created because no notification loop is running")
	})

	m.Run()

	t.Run("when notification loop is running subscription should be created", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, ch)
		require.True(t, ok)

		readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL1, testURL2)})
		m.Notify([]adminapi.DiscoveredAdminAPI{
			testDiscoveredAdminAPI(testURL1),
			testDiscoveredAdminAPI(testURL2),
		})

		select {
		case <-ch:
			require.Len(t, m.GatewayClients(), 2, "expected to get 2 clients after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
	})

	t.Run("when multiple subscriptions are created, each of them should receive notifications", func(t *testing.T) {
		sub1, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, sub1)
		require.True(t, ok)

		sub2, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, sub2)
		require.True(t, ok)

		readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedPending: intoTurnedPending(testURL2)})
		m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})

		select {
		case <-sub2:
			require.Len(t, m.GatewayClients(), 1, "expected to get 1 client after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
		select {
		case <-sub1:
			require.Len(t, m.GatewayClients(), 1, "expected to get 1 client after the update")
		case <-time.After(time.Second):
			t.Error("did not receive notification after gateway clients changes")
		}
	})

	t.Run("when the context gets cancelled, subscriber channel gets closed", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.NotNil(t, ch)
		require.True(t, ok)

		cancel()

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Error("subscription channel wasn't closed after cancelling the context")
		}
	})

	t.Run("when the context is cancelled, subscriptions cannot be created", func(t *testing.T) {
		ch, ok := m.SubscribeToGatewayClientsChanges()
		require.Nil(t, ch)
		require.False(t, ok)
	})
}

func TestAdminAPIClientsManager_ConcurrentNotify(t *testing.T) {
	readinessChecker := &mockReadinessChecker{}
	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL1)})
	testClient, err := adminapi.NewTestClient(testURL1)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := clients.NewAdminAPIClientsManager(ctx, zapr.NewLogger(zap.NewNop()), []*adminapi.Client{testClient}, readinessChecker)
	require.NoError(t, err)
	m.Run()

	// Run a goroutine that will call GatewayClients() every millisecond.
	go func() {
		for {
			select {
			case <-time.Tick(time.Millisecond):
				_ = m.GatewayClients()
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			go m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})
		}
	}()

	require.Eventually(t, func() bool {
		return readinessChecker.CallsCount() == 100
	}, time.Second, time.Millisecond)
}

func TestAdminAPIClientsManager_GatewayClientsChanges(t *testing.T) {
	testClient, err := adminapi.NewTestClient(testURL1)
	require.NoError(t, err)

	readinessChecker := &mockReadinessChecker{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := clients.NewAdminAPIClientsManager(ctx, zapr.NewLogger(zap.NewNop()), []*adminapi.Client{testClient}, readinessChecker)
	require.NoError(t, err)

	m.Run()
	<-m.Running()

	var receivedNotificationsCount atomic.Uint32
	ch, ok := m.SubscribeToGatewayClientsChanges()
	require.NotNil(t, ch)
	require.True(t, ok)

	// Run subscriber worker in a separate goroutine to consume notifications.
	go func() {
		for {
			select {
			case <-ch:
				receivedNotificationsCount.Add(1)
			case <-ctx.Done():
				return
			}
		}
	}()

	firstClientsSet := []adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)}
	secondClientsSet := []adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL2)}
	notificationsCountEventuallyEquals := func(expectedCount int) {
		require.Eventually(t, func() bool {
			if count := receivedNotificationsCount.Load(); count != uint32(expectedCount) {
				t.Logf("Received %d notifications, expected %d, waiting...", count, expectedCount)
				return false
			}
			return true
		}, time.Second, time.Millisecond, "expected to receive %d notifications", expectedCount)
	}
	requireLastReadinessCheckCall := func(expected readinessCheckCall) {
		call, ok := readinessChecker.LastCall()
		require.True(t, ok, "expected call to readiness checker")
		require.Equal(t, expected, call)
	}

	// Notify the first set of clients and make sure that the subscriber doesn't get notified as it was initial state.
	m.Notify(firstClientsSet)
	notificationsCountEventuallyEquals(0)
	require.Equal(t, 1, readinessChecker.CallsCount(), "expected readiness check on non-empty set of clients")
	requireLastReadinessCheckCall(readinessCheckCall{
		AlreadyCreatedURLs: []string{testURL1},
		PendingURLs:        []string{},
	})

	// Notify an empty set of clients and make sure that the subscriber get notified.
	m.Notify(nil)
	notificationsCountEventuallyEquals(1)
	require.Equal(t, 1, readinessChecker.CallsCount(), "no readiness check should be performed when notifying an empty set")

	// Notify an empty set again and make sure that the subscriber doesn't get notified as the state didn't change.
	m.Notify(nil)
	notificationsCountEventuallyEquals(1)
	require.Equal(t, 1, readinessChecker.CallsCount(), "no readiness check should be performed when notifying an empty set")

	// Notify the second set of clients without making the new one ready and make sure that the subscriber gets no notification.
	m.Notify(secondClientsSet)
	notificationsCountEventuallyEquals(1)
	requireLastReadinessCheckCall(readinessCheckCall{
		AlreadyCreatedURLs: []string{},
		PendingURLs:        []string{testURL2},
	})
	require.Equal(t, 2, readinessChecker.CallsCount(), "expected readiness check on non-empty set of clients")

	// Notify the second set of clients and make sure that the subscriber gets notified after the new one becomes ready.
	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL2)})
	m.Notify(secondClientsSet)
	notificationsCountEventuallyEquals(2)
	requireLastReadinessCheckCall(readinessCheckCall{
		AlreadyCreatedURLs: []string{},
		PendingURLs:        []string{testURL2},
	})
	require.Equal(t, 3, readinessChecker.CallsCount(), "expected readiness check on non-empty set of clients")

	m.Notify([]adminapi.DiscoveredAdminAPI{firstClientsSet[0], secondClientsSet[0]})
	notificationsCountEventuallyEquals(3)
	requireLastReadinessCheckCall(readinessCheckCall{
		AlreadyCreatedURLs: []string{testURL2},
		PendingURLs:        []string{testURL1},
	})
	require.Equal(t, 4, readinessChecker.CallsCount(), "expected readiness check on non-empty set of clients")
}

func TestAdminAPIClientsManager_PeriodicReadinessReconciliation(t *testing.T) {
	testClient, err := adminapi.NewTestClient(testURL1)
	require.NoError(t, err)

	readinessTicker := mocks.NewTicker()
	readinessChecker := &mockReadinessChecker{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := clients.NewAdminAPIClientsManager(
		ctx,
		zapr.NewLogger(zap.NewNop()),
		[]*adminapi.Client{testClient},
		readinessChecker,
		clients.WithReadinessReconciliationTicker(readinessTicker),
	)
	require.NoError(t, err)
	m.Run()
	<-m.Running()

	readinessCheckCallEventuallyMatches := func(expected readinessCheckCall) {
		require.Eventually(t, func() bool {
			lastCall, wasCalledAtAll := readinessChecker.LastCall()
			if !wasCalledAtAll {
				t.Log("Readiness checker was not called yet, waiting...")
				return false
			}
			if diff := cmp.Diff(expected, lastCall); diff != "" {
				t.Logf("Readiness checker was called with unexpected arguments: %s", diff)
				return false
			}
			return true
		}, time.Second, time.Millisecond)
	}

	// Trigger the first readiness check.
	readinessTicker.Add(clients.DefaultReadinessReconciliationInterval)
	readinessCheckCallEventuallyMatches(readinessCheckCall{
		AlreadyCreatedURLs: []string{testURL1},
		PendingURLs:        []string{},
	})
	require.Equal(t, 1, readinessChecker.CallsCount())

	// Notify with a new client and check the readiness check call was made as expected.
	m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1), testDiscoveredAdminAPI(testURL2)})
	readinessCheckCallEventuallyMatches(readinessCheckCall{
		AlreadyCreatedURLs: []string{testURL1},
		PendingURLs:        []string{testURL2},
	})
	require.Equal(t, 2, readinessChecker.CallsCount())

	// Trigger a next readiness check which will make testURL2 ready.
	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL2)})
	readinessTicker.Add(clients.DefaultReadinessReconciliationInterval)
	readinessCheckCallEventuallyMatches(readinessCheckCall{
		AlreadyCreatedURLs: []string{testURL1},
		PendingURLs:        []string{testURL2},
	})
	require.Equal(t, 3, readinessChecker.CallsCount())
	require.True(t, lo.ContainsBy(m.GatewayClients(), func(c *adminapi.Client) bool {
		return c.BaseRootURL() == testURL2
	}), "expected to find the new client in the manager's clients list after it became ready")
}

func testDiscoveredAdminAPI(address string) adminapi.DiscoveredAdminAPI {
	return adminapi.DiscoveredAdminAPI{
		Address: address,
		PodRef:  k8stypes.NamespacedName{Name: "pod-1", Namespace: "ns"},
	}
}
