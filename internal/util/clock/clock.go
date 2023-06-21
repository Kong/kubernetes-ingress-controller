package clock

import "time"

type System struct{}

func (System) Now() time.Time { return time.Now() }
