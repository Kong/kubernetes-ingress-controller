package dataplane

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
)

const testSynchronizerTick = time.Millisecond * 10

func TestSynchronizer(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Log("setting up a fake dataplane client to test the synchronizer")
	c := &fakeDataplaneClient{dbmode: dpconf.DBModePostgres}

	t.Log("configuring the dataplane synchronizer")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("initializing the dataplane synchronizer")
	sync, err := NewSynchronizer(
		zapr.NewLogger(zap.NewNop()),
		c,
		WithStagger(testSynchronizerTick),
		WithInitCacheSyncDuration(testSynchronizerTick),
	)
	require.NoError(t, err)
	assert.NotNil(t, sync)

	t.Log("verifying that a non-started dataplane synchronizer reports as not running")
	assert.False(t, sync.IsRunning())

	t.Log("verifying that postgres dp makes the synchronizer immediately ready")
	sync.dbMode = "postgres"
	assert.True(t, sync.IsReady())

	t.Log("verifying dbless mode means the synchronizer wont be ready until a config has been applied")
	sync.dbMode = "off"
	assert.False(t, sync.IsReady())

	t.Log("starting the dataplane synchronizer server")
	assert.NoError(t, sync.Start(ctx))
	assert.Eventually(t, func() bool { return sync.IsRunning() }, time.Second, testSynchronizerTick)
	assert.True(t, sync.NeedLeaderElection())

	t.Log("verifying that trying to start the dataplane synchronizer while it's already started fails")
	err = sync.Start(ctx)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "server is already running")

	t.Log("verifying that eventually the synchronizer reports as ready for a dbless dataplane")
	assert.Eventually(t, func() bool { return sync.IsReady() }, testSynchronizerTick*3, testSynchronizerTick)

	t.Log("verifying that the dataplane eventually receieves several successful updates from the synchronizer")
	assert.Eventually(t, func() bool {
		return c.totalUpdates() >= 5
	}, 10*testSynchronizerTick, testSynchronizerTick, "got %d updates, expected 5 or more", c.totalUpdates())

	t.Log("verifying that the server shuts down when the context is cancelled")
	cancel()
	assert.Eventually(t, func() bool { return !sync.IsRunning() }, time.Second, testSynchronizerTick)
	assert.Eventually(t, func() bool { return !sync.IsReady() }, time.Second, testSynchronizerTick)
	totalUpdatesSeenSoFar := c.totalUpdates()

	t.Log("verifying that the server can be started back up with a new context")
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	assert.NoError(t, sync.Start(ctx))
	assert.Eventually(t, func() bool { return sync.IsRunning() }, time.Second, testSynchronizerTick)
	assert.Eventually(t, func() bool { return sync.IsReady() }, time.Second, testSynchronizerTick)

	t.Log("verifying that a server that was restarted continues to send successful updates to the dataplane")
	assert.Eventually(t, func() bool { return c.totalUpdates() >= totalUpdatesSeenSoFar+3 }, time.Second, testSynchronizerTick)

	t.Log("verifying that the server can be shut down a second time")
	cancel()
	assert.Eventually(t, func() bool { return !sync.IsRunning() }, time.Second, testSynchronizerTick)
	assert.Eventually(t, func() bool { return !sync.IsReady() }, time.Second, testSynchronizerTick)
}

func TestSynchronizer_IsReadyDoesntBlockWhenDataPlaneIsBlocked(t *testing.T) {
	for _, dbMode := range []dpconf.DBMode{
		dpconf.DBModeOff,
		dpconf.DBModePostgres,
	} {
		dbMode := dbMode
		t.Run(fmt.Sprintf("dbmode=%s", dbMode), func(t *testing.T) {
			c := &fakeDataplaneClient{dbmode: dbMode, t: t}
			l, err := zap.NewDevelopment()
			require.NoError(t, err)
			s, err := NewSynchronizer(
				zapr.NewLogger(l),
				c,
				WithStagger(testSynchronizerTick),
				WithInitCacheSyncDuration(testSynchronizerTick),
			)
			require.NoError(t, err)
			require.NotNil(t, s)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			require.NoError(t, s.Start(ctx))

			t.Log("verifying the first update happened and the synchronizer is ready")
			require.Eventually(t, func() bool { return s.IsReady() }, testSynchronizerTick*10, testSynchronizerTick)

			const clientBlockTime = time.Second * 10

			t.Log("making the data plane calls block for significantly longer than the synchronizer tick")
			c.clientShouldBlockFor(clientBlockTime)
			updateCount := c.totalUpdates()

			t.Log("waiting for a blocking update to happen")
			require.Eventually(t, func() bool { return c.totalUpdates() > updateCount }, testSynchronizerTick*10, testSynchronizerTick)

			t.Log("verifying that IsReady is not blocking even though the client is blocked")
			// NOTE: Allow a little extra time for the synchronizer to return from IsReady() as
			// time duration on the magnitude of milliseconds can be flaky.
			// As long as it responds in a reasonable amount of time ( less than half the blocking time), we're good.
			require.Eventually(t, func() bool { return s.IsReady() }, clientBlockTime/2, testSynchronizerTick)
		})
	}
}

// fakeDataplaneClient fakes the dataplane.Client interface so that we can
// unit test the dataplane.Synchronizer.
type fakeDataplaneClient struct {
	dbmode                  dpconf.DBMode
	updateCount             atomic.Uint64
	lock                    sync.RWMutex
	clientCallBlockDuration time.Duration
	t                       *testing.T
}

func (c *fakeDataplaneClient) DBMode() dpconf.DBMode {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.clientCallBlockDuration > 0 {
		c.t.Logf("DBMode() blocking for %s", c.clientCallBlockDuration)
		time.Sleep(c.clientCallBlockDuration)
	}
	return c.dbmode
}

func (c *fakeDataplaneClient) Update(ctx context.Context) error {
	c.updateCount.Add(1)
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.clientCallBlockDuration > 0 {
		ch := time.After(c.clientCallBlockDuration)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.t.Logf("Update() blocking for %s", c.clientCallBlockDuration)
		}

		select {

		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
		}
	}
	return nil
}

func (c *fakeDataplaneClient) Shutdown(_ context.Context) error {
	return nil
}

func (c *fakeDataplaneClient) clientShouldBlockFor(d time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clientCallBlockDuration = d
}

func (c *fakeDataplaneClient) totalUpdates() int {
	return int(c.updateCount.Load())
}
