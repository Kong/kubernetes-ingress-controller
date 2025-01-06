package mocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTicker(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ticker := NewTicker()
		ch := ticker.Channel()
		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}

		ticker.Reset(time.Hour)

		t.Log("adding second should not tick when ticker has an interval of 1 hour")
		ticker.Add(time.Second)
		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}

		t.Log("adding 40 minutes should not tick when ticker has an interval of 1 hour")
		ticker.Add(40 * time.Minute)
		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}

		t.Log("adding 40 minutes should tick when 40 minutes already passed and ticker has an interval of 1 hour")
		ticker.Add(40 * time.Minute)
		select {
		case <-ch:
		case <-time.After(time.Second):
			require.FailNow(t, "expected a tick to happen but it didn't")
		}
	})

	t.Run("Reset", func(t *testing.T) {
		ticker := NewTicker()
		ch := ticker.Channel()

		t.Log("reseting ticker to 3 hour interval")
		ticker.Reset(3 * time.Hour)
		t.Log("adding second should not tick when ticker has an interval of 3 hours")
		ticker.Add(time.Second)
		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}

		t.Log("adding an hour should not tick when ticker has an interval of 3 hour")
		ticker.Add(time.Hour)
		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}

		t.Log("adding 2 hours should tick when ticker has an interval of 3 hour")
		ticker.Add(2 * time.Hour)
		select {
		case <-ch:
		case <-time.After(time.Second):
			require.FailNow(t, "expected a tick to happen but it didn't")
		}
	})

	t.Run("stop", func(t *testing.T) {
		ticker := NewTicker()
		ch := ticker.Channel()
		ticker.Stop()

		select {
		case <-ch:
			require.FailNow(t, "unexpected tick")
		default:
		}
	})
}
