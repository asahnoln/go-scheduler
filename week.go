package scheduler

import "time"

// Week holds schedules for all days in a week
type Week struct {
	Days [7]Schedule
}

// NewWeek returns new Week
func NewWeek() Week {
	return Week{}
}

// Add adds a Schedule to the Week for a given day
func (w Week) Add(d time.Weekday, s Schedule) Week {
	w.Days[d] = w.Days[d].AddSchedule(s)
	return w
}

// Day returns Schedule for given day in a week
// If d is not a valid Weekday, an empty Schedule is returned
func (w Week) Day(d time.Weekday) Schedule {
	return w.Days[d]
}
