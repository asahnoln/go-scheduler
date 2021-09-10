package scheduler

import (
	"fmt"
	"time"
)

// Range is a busy range with start and end times
type Range struct {
	Start, End time.Time
}

// NewRangeFromStrings creates Range with given string times
func NewRangeFromStrings(start, end string) (Range, error) {
	startTime, endTime, err := parseTimes(start, end)

	// This kinda overrides above error checking when testing error path.
	// Potential problem?
	if !startTime.Before(endTime) {
		return Range{}, fmt.Errorf("scheduler: want new range start time less than end time, got AddRange(%q, %q)", start, end)
	}

	return Range{startTime, endTime}, err
}

// StartString returns start time in string format
func (r Range) StartString() string {
	return r.Start.Format("15:04")
}

// EndString returns end time in string format
func (r Range) EndString() string {
	return r.End.Format("15:04")
}

func parseTimes(start, end string) (time.Time, time.Time, error) {
	var (
		startTime, endTime time.Time
		err                error
	)

	startTime, err = time.Parse(LayoutTime, start)
	if err != nil {
		return startTime, endTime, fmt.Errorf("scheduler: error parsing start time: %w", err)
	}

	endTime, err = time.Parse(LayoutTime, end)
	if err != nil {
		return startTime, endTime, fmt.Errorf("scheduler: error parsing end time: %w", err)
	}

	return startTime, endTime, nil
}
