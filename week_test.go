package scheduler_test

import (
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler"
	"github.com/asahnoln/go-scheduler/test"
)

func TestCreateWeekSchedule(t *testing.T) {
	w := scheduler.NewWeek()
	s, _ := scheduler.NewSchedule().Add("09:00", "14:00")

	w = w.Add(time.Monday, s)
	d := w.Day(time.Monday)

	test.AssertSameLength(t, 1, len(d))

	s, _ = scheduler.NewSchedule().Add("15:00", "18:00")
	w = w.Add(time.Monday, s)
	d = w.Day(time.Monday)

	test.AssertSameLength(t, 2, len(d))
}

// Kinda cumbersome test
func TestWholeWeek(t *testing.T) {
	tests := [7]struct {
		day       time.Weekday
		times     [][2]string
		wantTimes [][2]string
	}{
		{time.Monday,
			[][2]string{
				{"07:05", "08:15"},
				{"07:30", "09:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"07:05", "09:05"},
				{"18:20", "19:10"},
			}},
		{time.Tuesday,
			[][2]string{
				{"08:05", "08:45"},
				{"08:30", "09:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"08:05", "09:05"},
				{"18:20", "19:10"},
			}},
		{time.Wednesday,
			[][2]string{
				{"07:05", "07:15"},
				{"07:10", "08:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"07:05", "08:05"},
				{"18:20", "19:10"},
			}},
		{time.Thursday,
			[][2]string{
				{"06:05", "07:15"},
				{"06:30", "08:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"06:05", "08:05"},
				{"18:20", "19:10"},
			}},
		{time.Friday,
			[][2]string{
				{"08:05", "09:15"},
				{"08:30", "10:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"08:05", "10:05"},
				{"18:20", "19:10"},
			}},
		{time.Saturday,
			[][2]string{
				{"05:05", "08:15"},
				{"07:30", "09:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"05:05", "09:05"},
				{"18:20", "19:10"},
			},
		},
		{time.Sunday,
			[][2]string{
				{"04:05", "05:15"},
				{"04:30", "06:05"},
				{"18:20", "19:10"},
			},
			[][2]string{
				{"04:05", "06:05"},
				{"18:20", "19:10"},
			},
		},
	}

	// Assert that time is merging
	w := scheduler.NewWeek()
	for _, tt := range tests {
		for _, ts := range tt.times {
			s, _ := scheduler.NewSchedule().Add(ts[0], ts[1])
			w = w.Add(tt.day, s)
		}
	}

	for _, tt := range tests {
		t.Run(tt.day.String(), func(t *testing.T) {
			d := w.Day(tt.day)

			test.AssertSameLength(t, len(tt.wantTimes), len(d))
			for i, want := range tt.wantTimes {
				test.AssertSameString(t, want[0], d[i].StartString(), "want start time %q, got %q")
				test.AssertSameString(t, want[1], d[i].EndString(), "want end time %q, got %q")
			}
		})
	}
}

// Seems like map or switch doesn't make a difference for some reason
// Slice is kinda faster?
func BenchmarkWeek(b *testing.B) {
	days := [7]time.Weekday{
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
		time.Sunday,
	}
	for i := 0; i < b.N; i++ {
		w := scheduler.NewWeek()
		for _, d := range days {
			s, _ := scheduler.NewSchedule().Add("09:00", "10:00")
			w := w.Add(d, s)
			w.Day(d)
		}
	}
}
