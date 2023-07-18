package clients_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/clients"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"golang.org/x/exp/slices"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
)

type mockReadinessChecker struct {
	nextResult clients.ReadinessCheckResult
	lock       sync.RWMutex
}

func (m *mockReadinessChecker) CheckReadiness(
	context.Context,
	[]clients.AlreadyCreatedClient,
	[]adminapi.DiscoveredAdminAPI,
) clients.ReadinessCheckResult {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.nextResult
}

func (m *mockReadinessChecker) LetChecksReturn(result clients.ReadinessCheckResult) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.nextResult = result
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

	logger := logrus.New()
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
	_, err := clients.NewAdminAPIClientsManager(context.Background(), logrus.New(), nil, &mockReadinessChecker{})
	require.ErrorContains(t, err, "at least one initial client must be provided")
}

func TestAdminAPIClientsManager_NotRunningNotifyLoop(t *testing.T) {
	t.Parallel()

	testClient, err := adminapi.NewTestClient("localhost:8080")
	require.NoError(t, err)
	m, err := clients.NewAdminAPIClientsManager(
		context.Background(),
		logrus.New(),
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
		logrus.New(),
		[]*adminapi.Client{testClient},
		&mockReadinessChecker{},
	)
	require.NoError(t, err)
	require.Len(t, m.GatewayClients(), 1, "expecting one initial client")
	require.Equal(t, m.GatewayClientsCount(), 1, "expecting one initial client")

	konnectTestClient := &adminapi.KonnectClient{}
	m.SetKonnectClient(konnectTestClient)
	require.Len(t, m.GatewayClients(), 1, "konnect client should not be returned from GatewayClients")
	require.Equal(t, m.GatewayClientsCount(), 1, "konnect client should not be counted in GatewayClientsCount")
	require.Equal(t, konnectTestClient, m.KonnectClient(), "konnect client should be returned from KonnectClient")
}

func TestAdminAPIClientsManager_SubscribeToGatewayClientsChanges(t *testing.T) {
	t.Parallel()

	readinessChecker := &mockReadinessChecker{}
	testClient, err := adminapi.NewTestClient("http://localhost:8080")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	m, err := clients.NewAdminAPIClientsManager(ctx, logrus.New(), []*adminapi.Client{testClient}, readinessChecker)
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
	t.Parallel()

	readinessChecker := &mockReadinessChecker{}
	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL1)})
	testClient, err := adminapi.NewTestClient(testURL1)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := clients.NewAdminAPIClientsManager(ctx, logrus.New(), []*adminapi.Client{testClient}, readinessChecker)
	require.NoError(t, err)
	m.Run()

	var receivedNotificationsCount atomic.Uint32
	ch, ok := m.SubscribeToGatewayClientsChanges()
	require.NotNil(t, ch)
	require.True(t, ok)

	// Run subscriber worker in a separate goroutine to consume notifications.
	go func() {
		for {
			select {
			case <-ch:
				// Call GatewayClients() here to make sure that we can access the clients safely
				// from the subscriber goroutine without causing a deadlock in the notify loop.
				_ = m.GatewayClients()
				receivedNotificationsCount.Add(1)
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			// Swap between ready and pending interchangeably depending on the iteration to trigger a change.
			if pickEven := i%2 == 0; pickEven {
				readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL1)})
			} else {
				readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedPending: intoTurnedPending(testURL1)})
			}
			m.Notify([]adminapi.DiscoveredAdminAPI{testDiscoveredAdminAPI(testURL1)})
		}()
	}

	require.Eventually(t, func() bool {
		if count := receivedNotificationsCount.Load(); count < 2 {
			t.Logf("Received %d notifications, expected at least 2, waiting...", count)
			return false
		}
		return true
	}, time.Second, time.Millisecond, "expected to receive at least 2 notifications")
}

func TestAdminAPIClientsManager_NotifiesSubscribersOnlyWhenGatewayClientsChange(t *testing.T) {
	testClient, err := adminapi.NewTestClient(testURL1)
	require.NoError(t, err)

	readinessChecker := &mockReadinessChecker{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m, err := clients.NewAdminAPIClientsManager(ctx, logrus.New(), []*adminapi.Client{testClient}, readinessChecker)
	require.NoError(t, err)
	m.Run()

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

	firstClientsSet := []adminapi.DiscoveredAdminAPI{
		testDiscoveredAdminAPI(testURL1),
	}
	secondClientsSet := []adminapi.DiscoveredAdminAPI{
		testDiscoveredAdminAPI(testURL2),
	}
	notificationsCountEventuallyEquals := func(expectedCount int) {
		require.Eventually(t, func() bool {
			if count := receivedNotificationsCount.Load(); count != uint32(expectedCount) {
				t.Logf("Received %d notifications, expected %d, waiting...", count, expectedCount)
				return false
			}
			return true
		}, time.Second, time.Millisecond, "expected to receive %d notifications", expectedCount)
	}

	// Notify the first set of clients and make sure that the subscriber doesn't get notified as it was initial state.
	m.Notify(firstClientsSet)
	notificationsCountEventuallyEquals(0)

	// Notify an empty set of clients and make sure that the subscriber get notified.
	m.Notify(nil)
	notificationsCountEventuallyEquals(1)

	// Notify an empty set again and make sure that the subscriber doesn't get notified as the state didn't change.
	m.Notify(nil)
	notificationsCountEventuallyEquals(1)

	// Notify the second set of clients without making the new one ready and make sure that the subscriber gets no notification.
	m.Notify(secondClientsSet)
	notificationsCountEventuallyEquals(1)

	// Notify the second set of clients and make sure that the subscriber gets notified after the new one becomes ready.
	readinessChecker.LetChecksReturn(clients.ReadinessCheckResult{ClientsTurnedReady: intoTurnedReady(testURL2)})
	m.Notify(secondClientsSet)
	notificationsCountEventuallyEquals(2)
}

func testDiscoveredAdminAPI(address string) adminapi.DiscoveredAdminAPI {
	return adminapi.DiscoveredAdminAPI{
		Address: address,
		PodRef:  k8stypes.NamespacedName{Name: "pod-1", Namespace: "ns"},
	}
}
