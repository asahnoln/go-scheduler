package scheduler

import "time"

// Range is a busy range with start and end times
type Range struct {
	Start, End time.Time
}

// StartString returns start time in string format
func (r Range) StartString() string {
	return r.Start.Format("15:04")
}

// EndString returns end time in string format
func (r Range) EndString() string {
	return r.End.Format("15:04")
}
