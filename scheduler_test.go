package scheduler_test

import (
	"testing"

	"github.com/asahnoln/go-scheduler"
)

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name      string
		times     [][2]string
		wantTimes [][2]string
	}{
		{
			"in order, separate times",
			[][2]string{
				{"09:00", "12:00"},
				{"18:00", "20:00"},
			},
			[][2]string{
				{"09:00", "12:00"},
				{"18:00", "20:00"},
			},
		},
		{
			"added in order, overlapping ending and beginning",
			[][2]string{
				{"09:00", "12:00"},
				{"11:00", "14:00"},
			},
			[][2]string{
				{"09:00", "14:00"},
			},
		},
		{
			"reverse order, overlapping beginning and ending",
			[][2]string{
				{"11:00", "14:00"},
				{"09:00", "12:00"},
			},
			[][2]string{
				{"09:00", "14:00"},
			},
		},
		{
			"in order, touching ending and beginning",
			[][2]string{
				{"09:00", "12:00"},
				{"12:00", "14:00"},
			},
			[][2]string{
				{"09:00", "14:00"},
			},
		},
		{
			"reverse order, touching beginning and ending",
			[][2]string{
				{"12:00", "14:00"},
				{"09:00", "12:00"},
			},
			[][2]string{
				{"09:00", "14:00"},
			},
		},
		{
			"separate time, added in order, overlapping",
			[][2]string{
				{"09:00", "12:00"},
				{"15:00", "18:00"},
				{"16:00", "19:00"},
			},
			[][2]string{
				{"09:00", "12:00"},
				{"15:00", "19:00"},
			},
		},
		{
			"separate time, added not in order, overlapping and touching",
			[][2]string{
				{"20:00", "22:00"}, // touches 22:00
				{"15:00", "18:00"}, // overlaps 16:00
				{"09:00", "12:00"}, // separate
				{"16:00", "19:00"}, // overlaps 18:00
				{"14:00", "18:00"}, // overlaps 15:00
				{"22:00", "23:00"}, // touches 22:00
			},
			[][2]string{
				{"09:00", "12:00"},
				{"14:00", "19:00"},
				{"20:00", "23:00"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scheduler.NewSchedule()

			for _, args := range tt.times {
				s, _ = s.AddRange(args[0], args[1])
			}
			t.Logf("times added: %v", tt.times)
			for _, r := range s {
				t.Logf("times in Schedule: %v - %v", r.Start(), r.End())
			}

			assertSameLength(t, len(tt.wantTimes), len(s))

			for i, ts := range tt.wantTimes {
				r := s[i]
				assertSameString(t, ts[0], r.Start(), "want range start %q, got %q")
				assertSameString(t, ts[1], r.End(), "want range end %q, got %q")
			}
		})
	}
}

// TODO: Assert Error objects
func TestRangeStartLessThanEnd(t *testing.T) {
	_, err := scheduler.NewSchedule().AddRange("14:00", "09:00")
	assertError(t, err, "expect error, because given start time is greater than end time, got %v")

	_, err = scheduler.NewSchedule().AddRange("14:00", "14:00")
	assertError(t, err, "expect error, because given start time is equal to end time, got %v")
}

// TODO: Think on optimization
func BenchmarkMerging(b *testing.B) {
	s := scheduler.NewSchedule()
	for i := 0; i < b.N; i++ {
		s, _ = s.AddRange("09:00", "14:00")
	}
}

func assertError(t testing.TB, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatalf(message, err)
	}
}

func assertNoError(t testing.TB, err error, message string) {
	t.Helper()

	if err != nil {
		t.Fatalf(message, err)
	}
}

func assertSameLength(t testing.TB, want, got int) {
	t.Helper()

	if want != got {
		t.Fatalf("want range count %d, got %d", want, got)
	}
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
