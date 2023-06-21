package clock

import "time"

func NewTicker() *TimeTicker {
	return &TimeTicker{
		ticker: time.NewTicker(1000 * time.Hour),
	}
}

type TimeTicker struct {
	ticker *time.Ticker
}

func (t *TimeTicker) Stop() {
	t.ticker.Stop()
}

func (t *TimeTicker) Channel() <-chan time.Time {
	return t.ticker.C
}

func (t *TimeTicker) Reset(d time.Duration) {
	t.ticker.Reset(d)
}
