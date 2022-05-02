package dataplane

import (
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSynchronizer(t *testing.T) {
	t.Log("setting up a fake dataplane client to test the synchronizer")
	c := &fakeDataplaneClient{dbmode: "postgres"}

	t.Log("configuring the dataplane synchronizer")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("initializing the dataplane synchronizer")
	stagger := time.Millisecond * 200
	sync, err := NewSynchronizerWithStagger(logrus.New(), c, stagger)
	assert.NoError(t, err)
	assert.NotNil(t, sync)

	t.Log("verifying that a non-started dataplane synchronizer reports as not running")
	assert.False(t, sync.IsRunning())

	t.Log("verifying that postgres dp makes the synchronizer immediately ready")
	c.dbmode = "postgres"
	assert.True(t, sync.IsReady())

	t.Log("verifying dbless mode means the synchronizer wont be ready until a config has been applied")
	c.dbmode = "off"
	assert.False(t, sync.IsReady())

	t.Log("starting the dataplane synchronizer server")
	assert.NoError(t, sync.Start(ctx))
	assert.Eventually(t, func() bool { return sync.IsRunning() }, time.Second, time.Millisecond*200)
	assert.True(t, sync.NeedLeaderElection())

	t.Log("verifying that trying to start the dataplane synchronizer while it's already started fails")
	err = sync.Start(ctx)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "server is already running")

	t.Log("verifying that eventually the synchronizer reports as ready for a dbless dataplane")
	assert.Eventually(t, func() bool { return sync.IsReady() }, stagger*2, time.Millisecond*200)

	t.Log("verifying that the dataplane eventually receieves several successful updates from the synchronizer")
	assert.Eventually(t, func() bool {
		return c.totalUpdates() >= 5
	}, stagger*2, time.Millisecond*200)

	t.Log("verifying that the server shuts down when the context is cancelled")
	cancel()
	assert.Eventually(t, func() bool { return !sync.IsRunning() }, time.Second, time.Millisecond*200)
	assert.Eventually(t, func() bool { return !sync.IsReady() }, time.Second, time.Millisecond*200)
	totalUpdatesSeenSoFar := c.totalUpdates()

	t.Log("verifying that the server can be started back up with a new context")
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	assert.NoError(t, sync.Start(ctx))
	assert.Eventually(t, func() bool { return sync.IsRunning() }, time.Second, time.Millisecond*200)
	assert.Eventually(t, func() bool { return sync.IsReady() }, time.Second, time.Millisecond*200)

	t.Log("verifying that a server that was restarted continues to send successful updates to the dataplane")
	assert.Eventually(t, func() bool { return c.totalUpdates() >= totalUpdatesSeenSoFar+3 }, time.Second, time.Millisecond*200)

	t.Log("verifying that the server can be shut down a second time")
	cancel()
	assert.Eventually(t, func() bool { return !sync.IsRunning() }, time.Second, time.Millisecond*200)
	assert.Eventually(t, func() bool { return !sync.IsReady() }, time.Second, time.Millisecond*200)
}

// fakeDataplaneClient fakes the dataplane.Client interface so that we can
// unit test the dataplane.Synchronizer.
type fakeDataplaneClient struct {
	dbmode      string
	updateCount int
	lock        sync.RWMutex
}

func (c *fakeDataplaneClient) DBMode() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.dbmode
}

func (c *fakeDataplaneClient) Update(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.updateCount++
	return nil
}

func (c *fakeDataplaneClient) totalUpdates() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.updateCount
}
