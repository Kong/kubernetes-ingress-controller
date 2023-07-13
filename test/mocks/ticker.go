package mocks

import (
	"sync"
	"time"
)

const (
	// This is irrelevant for the ticker, but we need to pass something to NewTicker.
	// The reason for this is that the ticker is used in the license agent, which
	// uses a non trivial logic to determine the polling period based on the state
	// of license retrieval.
	// This might be changed in the future if it doesn't fit the future needs.
	initialTickerDuration = 1000 * time.Hour
)

func NewTicker() *Ticker {
	now := time.Now()

	ticker := &Ticker{
		sigClose: make(chan struct{}),
		d:        initialTickerDuration,
		ch:       make(chan time.Time, 1),
		time:     now,
		lastTick: now,
	}

	return ticker
}

type Ticker struct {
	lock     sync.RWMutex
	sigClose chan struct{}
	d        time.Duration
	ch       chan time.Time
	time     time.Time
	lastTick time.Time
}

func (t *Ticker) Stop() {
	close(t.sigClose)
}

func (t *Ticker) Channel() <-chan time.Time {
	return t.ch
}

func (t *Ticker) Now() time.Time {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.time
}

func (t *Ticker) Reset(d time.Duration) {
	select {
	case <-t.sigClose:
		return
	default:
	}

	now := time.Now()

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastTick = now
	t.time = now
	t.d = d
}

func (t *Ticker) Add(d time.Duration) {
	select {
	case <-t.sigClose:
		return
	default:
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.time = t.time.Add(d)

	if t.time.Compare(t.lastTick.Add(t.d)) >= 0 {
		select {
		case <-t.sigClose:
		case t.ch <- t.time:
		}
		t.lastTick = t.time
	}
}
