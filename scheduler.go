package scheduler

import (
	"sort"
)

// LayoutTime is common time format for parsing given range strings to time.Time
const LayoutTime = "15:04"

// Schedule is a slice of ranges
type Schedule []Range

// NewSchedule returns new Schedule
func NewSchedule() Schedule {
	return Schedule{}
}

// Add adds a new Range with given start and end times to the Schedule.
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
func (s Schedule) Add(start, end string) (Schedule, error) {
	r, err := NewRangeFromStrings(start, end)
	if err != nil {
		return s, err
	}

	return s.AddRange(r), nil
}

// AddRange adds Range r to Schedule s and merges it if needed
func (s Schedule) AddRange(r Range) Schedule {
	s = s.merge(r)
	sort.Slice(s, func(i, j int) bool {
		return s[i].Start.Before(s[j].Start)
	})

	return s
}

// AddSchedule merges Schedule s with Schedule c
func (s Schedule) AddSchedule(c Schedule) Schedule {
	for _, r := range c {
		s = s.merge(r)
	}
	return s
}

func (s Schedule) merge(r Range) Schedule {
	newS := NewSchedule()
	if l := len(s); l > 0 {
		for _, p := range s {
			switch {
			case p.Start.Before(r.Start) && p.End.After(r.End):
				r.Start = p.Start
				r.End = p.End
			case p.Start.Before(r.Start) && (p.End.After(r.Start) || p.End.Equal((r.Start))):
				r.Start = p.Start
			case p.End.After(r.End) && (p.Start.Before(r.End) || p.Start.Equal(r.End)):
				r.End = p.End
			case p.Start.After(r.Start) && p.End.Before(r.End) || p.Start.Equal(r.Start) && p.End.Equal(r.End):
			default:
				newS = append(newS, p)
			}
		}
	}
	newS = append(newS, r)
	return newS
}
