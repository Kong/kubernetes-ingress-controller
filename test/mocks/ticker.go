package mocks

import (
	"sync"
	"time"
)

func NewTicker() *Ticker {
	now := time.Now()

	ticker := &Ticker{
		sigTime:  make(chan time.Time),
		sigClose: make(chan struct{}, 1),
		sigAdd:   make(chan time.Duration),
		sigReset: make(chan time.Duration),
		d:        1000 * time.Hour,
		ch:       make(chan time.Time, 1),
		time:     now,
		lastTick: now,
	}

	return ticker
}

type Ticker struct {
	lock     sync.RWMutex
	sigTime  chan time.Time
	sigClose chan struct{}
	sigReset chan time.Duration
	sigAdd   chan time.Duration
	d        time.Duration
	ch       chan time.Time
	time     time.Time
	lastTick time.Time
}

func (t *Ticker) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()

	close(t.sigTime)
	close(t.sigAdd)
	close(t.sigReset)
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
	now := time.Now()

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastTick = now
	t.time = now
	t.d = d
}

func (t *Ticker) Add(d time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.time = t.time.Add(d)

	if t.time.Compare(t.lastTick.Add(t.d)) >= 0 {
		t.ch <- t.time
		t.lastTick = t.time
	}
}
