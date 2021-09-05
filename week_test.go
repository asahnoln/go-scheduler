package scheduler_test

import (
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler"
)

func TestCreateWeekSchedule(t *testing.T) {
	w := scheduler.NewWeek()
	s, _ := scheduler.NewSchedule().Add("09:00", "14:00")

	w = w.Add(time.Monday, s)
	d := w.Day(time.Monday)

	assertSameLength(t, 1, len(d))

	s, _ = scheduler.NewSchedule().Add("15:00", "18:00")
	w = w.Add(time.Monday, s)
	d = w.Day(time.Monday)

	assertSameLength(t, 2, len(d))
}

func TestWholeWeek(t *testing.T) {
	tests := [7]struct {
		day   time.Weekday
		times [2]string
	}{
		{time.Monday,
			[2]string{"07:10", "08:20"}},
		{time.Tuesday,
			[2]string{"08:20", "09:30"}},
		{time.Wednesday,
			[2]string{"09:30", "10:40"}},
		{time.Thursday,
			[2]string{"10:50", "11:00"}},
		{time.Friday,
			[2]string{"12:10", "13:20"}},
		{time.Saturday,
			[2]string{"14:30", "15:40"}},
		{time.Sunday,
			[2]string{"16:50", "17:00"}},
	}

	w := scheduler.NewWeek()
	for _, tt := range tests {
		s, _ := scheduler.NewSchedule().Add(tt.times[0], tt.times[1])
		w = w.Add(tt.day, s)
	}

	for _, tt := range tests {
		t.Run(tt.day.String(), func(t *testing.T) {
			d := w.Day(tt.day)

			assertSameLength(t, 1, len(d))
			assertSameString(t, tt.times[0], d[0].StartString(), "want start time %q, got %q")
			assertSameString(t, tt.times[1], d[0].EndString(), "want end time %q, got %q")
		})
	}
}
