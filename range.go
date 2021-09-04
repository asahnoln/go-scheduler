package scheduler

import "time"

// Range is a busy range with start and end times
type Range struct {
	start, end time.Time
}

// Start returns start time in string format
func (r Range) Start() string {
	return r.start.Format("15:04")
}

// End returns end time in string format
func (r Range) End() string {
	return r.end.Format("15:04")
}
