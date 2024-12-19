package clock

import "time"

const (
	// This is irrelevant for the ticker, but we need to pass something to NewTicker.
	// The reason for this is that the ticker is used in the license agent, which
	// uses a non trivial logic to determine the polling period based on the state
	// of license retrieval.
	// This might be changed in the future if it doesn't fit the future needs.
	initialTickerDuration = 1000 * time.Hour
)

func NewTicker() *TimeTicker {
	return &TimeTicker{
		ticker: time.NewTicker(initialTickerDuration),
	}
}

func NewTickerWithDuration(d time.Duration) *TimeTicker {
	return &TimeTicker{
		ticker: time.NewTicker(d),
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
