package scheduler

import "time"

type Schedule []Range

type Range struct {
	start, end string
}

func NewSchedule() Schedule {
	return Schedule{}
}

func (s Schedule) AddRange(start, end string) (Schedule, error) {
	newRange := Range{start, end}
	var err error

	if l := len(s); l > 0 {
		for i, r := range s {
			rStart, err := time.Parse("15:04", r.start)
			if err != nil {
				return s, err
			}
			rEnd, err := time.Parse("15:04", r.end)
			if err != nil {
				return s, err
			}
			newStart, err := time.Parse("15:04", newRange.start)
			if err != nil {
				return s, err
			}
			newEnd, err := time.Parse("15:04", newRange.end)
			if err != nil {
				return s, err
			}

			// TODO: Fix conditions
			if rEnd.After(newStart) || rEnd.Equal(newStart) {
				s[i].end = newRange.end
			} else if rStart.Before(newEnd) || rStart.Equal(newEnd) {
				s[i].start = newRange.start
			} else {
				s = append(s, newRange)
			}
		}
	} else {
		s = append(s, newRange)
	}

	return s, err
}

func (r Range) Start() string {
	return r.start
}

func (r Range) End() string {
	return r.end
}
