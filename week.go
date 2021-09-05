package scheduler

import "time"

// Week holds schedules for all days in a week
type Week struct {
	mon, tue, wed, thu, fri, sat, sun Schedule
}

// NewWeek returns new Week
func NewWeek() Week {
	return Week{}
}

// Add adds a Schedule to the Week for a given day
func (w Week) Add(d time.Weekday, s Schedule) Week {
	switch d {
	case time.Monday:
		w.mon = w.mon.AddSchedule(s)
	case time.Tuesday:
		w.tue = w.tue.AddSchedule(s)
	case time.Wednesday:
		w.wed = w.wed.AddSchedule(s)
	case time.Thursday:
		w.thu = w.thu.AddSchedule(s)
	case time.Friday:
		w.fri = w.fri.AddSchedule(s)
	case time.Saturday:
		w.sat = w.sat.AddSchedule(s)
	case time.Sunday:
		w.sun = w.sun.AddSchedule(s)
	}

	return w
}

// Day returns Schedule for given day in a week
// If d is not a valid Weekday, an empty Schedule is returned
func (w Week) Day(d time.Weekday) Schedule {
	switch d {
	case time.Monday:
		return w.mon
	case time.Tuesday:
		return w.tue
	case time.Wednesday:
		return w.wed
	case time.Thursday:
		return w.thu
	case time.Friday:
		return w.fri
	case time.Saturday:
		return w.sat
	case time.Sunday:
		return w.sun
	}

	return Schedule{}
}
