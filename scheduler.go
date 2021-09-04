package scheduler

import (
	"fmt"
	"sort"
	"time"
)

type Schedule []Range

type Range struct {
	start, end time.Time
}

func NewSchedule() Schedule {
	return Schedule{}
}

func (s Schedule) AddRange(start, end string) (Schedule, error) {
	startTime, err := time.Parse("15:04", start)
	if err != nil {
		return s, err
	}
	endTime, err := time.Parse("15:04", end)
	if err != nil {
		return s, err
	}

	if !startTime.Before(endTime) {
		return s, fmt.Errorf("scheduler: want new range start time less than end time, got AddRange(%q, %q)", start, end)
	}

	newRange := Range{startTime, endTime}
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

	sort.Slice(newS, func(i, j int) bool {
		return newS[i].start.Before(newS[j].start)
	})
	return newS, err
}

func (r Range) Start() string {
	return r.start.Format("15:04")
}

func (r Range) End() string {
	return r.end.Format("15:04")
}
