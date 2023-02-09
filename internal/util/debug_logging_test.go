package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebugLoggerStiflesConsecutiveEntries(t *testing.T) {
	// initialize the debug logger with no backoff time, but a limit of 3 consecutive redudant entries
	buf := new(bytes.Buffer)
	log := MakeDebugLoggerWithReducedRedudancy(buf, &logrus.JSONFormatter{}, 3, time.Millisecond*0)
	assert.True(t, log.IsLevelEnabled(logrus.DebugLevel))
	assert.False(t, log.IsLevelEnabled(logrus.TraceLevel))

	// spam the logger with redudant entries and validate that only 3 entries (the limit) were actually emitted
	for i := 0; i < 100; i++ {
		log.Info("test")
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 3)

	// validate the logging data integrity
	for _, line := range lines {
		var entry map[string]string
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		assert.Equal(t, "test", entry["msg"])
	}
}

func TestDebugLoggerResetsConsecutiveEntries(t *testing.T) {
	// initialize the debug logger with no backoff time, but a limit of 5 consecutive redudant entries
	buf := new(bytes.Buffer)
	log := MakeDebugLoggerWithReducedRedudancy(buf, &logrus.JSONFormatter{}, 5, time.Millisecond*0)
	assert.True(t, log.IsLevelEnabled(logrus.DebugLevel))
	assert.False(t, log.IsLevelEnabled(logrus.TraceLevel))

	// spam the logger with redudant entries and validate that only 3 entries (the limit) were actually emitted
	for i := 0; i < 100; i++ {
		if i%5 == 0 {
			log.Info("break")
			continue
		}
		log.Info("test")
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 100)

	// validate the logging data integrity
	for i, line := range lines {
		var entry map[string]string
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		if i%5 == 0 {
			assert.Equal(t, "break", entry["msg"])
		} else {
			assert.Equal(t, "test", entry["msg"])
		}
	}
}

func TestDebugLoggerStiflesEntriesWhichAreTooFrequent(t *testing.T) {
	// initialize the debug logger with no consecutive entry backoff, but a time backoff of 30m
	buf := new(bytes.Buffer)
	log := MakeDebugLoggerWithReducedRedudancy(buf, &logrus.JSONFormatter{}, 0, time.Minute*30)

	// spam the logger, but validate that only one entry gets printed within the backoff timeframe
	for i := 0; i < 100; i++ {
		log.Debug("unique")
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 1)

	// validate the log entry
	var entry map[string]string
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &entry))
	assert.Equal(t, "unique", entry["msg"])
}

func TestDebugLoggerStopsStiflingEntriesAfterBackoffExpires(t *testing.T) {
	// setup backoffs and determine start/stop times
	start := time.Now()
	backoff := time.Millisecond * 100
	stop := start.Add(backoff)

	// initialize the debug logger with no consecutive entry backoff, but a time based backoff
	buf := new(bytes.Buffer)
	log := MakeDebugLoggerWithReducedRedudancy(buf, &logrus.JSONFormatter{}, 0, backoff)

	// spam the logger and validate that the testing environment didn't take longer than 100ms to process this
	for i := 0; i < 100; i++ {
		log.Debug("unique")
	}
	assert.True(t, time.Now().Before(stop),
		"validate that the resource contention in the testing environment is not overt")

	// verify that a new backoff period started and that all lines beyond the original were stifled
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 1)

	// validate the log data integrity
	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &entry))
	assert.Equal(t, "unique", entry["msg"])

	// wait until the backoff time is up and validate that it will allow one entry of the previous
	// redundant log entry to be emitted now that the backoff is over.
	time.Sleep(backoff)
	for i := 0; i < 1; i++ {
		log.Debug("second-unique")
	}

	// verify that a new backoff period started and that all lines beyond the original were stifled
	lines = strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 2)

	// validate the log data integrity
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &entry))
	assert.Equal(t, "unique", entry["msg"])
	require.NoError(t, json.Unmarshal([]byte(lines[1]), &entry))
	assert.Equal(t, "second-unique", entry["msg"])
}

func TestDebugLoggerThreadSafety(t *testing.T) {
	buf := &threadSafeBuffer{buf: new(bytes.Buffer), l: &sync.RWMutex{}}
	log := MakeDebugLoggerWithReducedRedudancy(buf, &logrus.JSONFormatter{}, 0, time.Minute*30)

	const total = 100
	writeToLogger := func() {
		// spam the logger concurrently across several goroutines to ensure no dataraces
		wg := &sync.WaitGroup{}
		wg.Add(total)
		for i := 0; i < total; i++ {
			go func() {
				defer wg.Done()
				log.Debug("unique")
			}()
		}
		wg.Wait()
	}
	assert.Eventually(t, func() bool {
		writeToLogger()

		if !strings.Contains(buf.String(), "unique") {
			return false
		}

		// Ensure that _some_ lines have been stifled. The actual number is not deterministic.
		lines := strings.Split(buf.String(), "\n")
		if !(len(lines) < total) {
			t.Logf("we haven't filtered any logs out, since the logger is indeterministic in nature let's retry")
			buf.buf.Reset()
			return false
		}
		return true
	}, time.Second, time.Millisecond)
}

// -----------------------------------------------------------------------------
// Private Types - Test Helpers
// -----------------------------------------------------------------------------

type threadSafeBuffer struct {
	buf *bytes.Buffer
	l   *sync.RWMutex
}

func (b *threadSafeBuffer) Read(p []byte) (n int, err error) {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.buf.Read(p)
}

func (b *threadSafeBuffer) Write(p []byte) (n int, err error) {
	b.l.Lock()
	defer b.l.Unlock()
	return b.buf.Write(p)
}

func (b *threadSafeBuffer) String() string {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.buf.String()
}
