package scheduler

import (
	"fmt"
	"sort"
	"time"
)

const LayoutTime = "15:04"

// Schedule is a slice of ranges
type Schedule []Range

// NewSchedule returns new Schedule
func NewSchedule() Schedule {
	return Schedule{}
}

// AddRange adds a new Range with given start and end times to the Schedule.
//
// Times parsed from strings to time.Time.
//
// New Range is merged with other Ranges in Schedule.
//
// If new Range overlaps with the ending or the beginning of other Range,
// new Range gets expanded accordingly as a sum of two ranges:
// 09:00-12:00 and 11:00-14:00 becomes 09:00-14:00 and so on.
//
// Schedule Ranges designed to work with clock time during one day,
// but they potentially might work with longer ranges, taking days, months or years.
func (s Schedule) AddRange(start, end string) (Schedule, error) {
	startTime, endTime, err := parseTimes(start, end)
	if err != nil {
		return s, err
	}

	if !startTime.Before(endTime) {
		return s, fmt.Errorf("scheduler: want new range start time less than end time, got AddRange(%q, %q)", start, end)
	}

	newS := s.merge(Range{startTime, endTime})
	sort.Slice(newS, func(i, j int) bool {
		return newS[i].start.Before(newS[j].start)
	})

	return newS, err
}

func (s Schedule) merge(newRange Range) Schedule {
	newS := NewSchedule()
	if l := len(s); l > 0 {
		for _, r := range s {
			switch {
			case r.start.Before(newRange.start) && (r.end.After(newRange.start) || r.end.Equal((newRange.start))):
				newRange.start = r.start
			case r.end.After(newRange.end) && (r.start.Before(newRange.end) || r.start.Equal(newRange.end)):
				newRange.end = r.end
			default:
				newS = append(newS, r)
			}
		}
	}
	newS = append(newS, newRange)
	return newS
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
