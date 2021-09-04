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
	newS := NewSchedule()
	newRange := Range{start, end}
	var err error

	if l := len(s); l > 0 {
		for _, r := range s {
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

			if rStart.Before(newStart) && rEnd.Before(newStart) || rStart.After(newEnd) && rEnd.After(newEnd) {
				newS = append(newS, r)
			} else if rStart.Before(newStart) && (rEnd.After(newStart) || rEnd.Equal((newStart))) {
				newRange.start = r.start
			} else if rEnd.After(newEnd) && (rStart.Before(newEnd) || rStart.Equal(newEnd)) {
				newRange.end = r.end
			}
		}
	}
	newS = append(newS, newRange)

	return newS, err
}

func (r Range) Start() string {
	return r.start
}

func (r Range) End() string {
	return r.end
}
