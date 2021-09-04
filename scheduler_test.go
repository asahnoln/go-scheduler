package scheduler_test

import (
	"testing"

	"github.com/asahnoln/go-scheduler"
)

func TestCreateSchedule(t *testing.T) {
	s := scheduler.NewSchedule()
	s, err := s.AddRange("09:00", "12:00")
	assertNoError(t, err, "unexpected error while adding range: %v")

	assertSameLength(t, 1, len(s))

	r := s[0]
	assertSameString(t, "09:00", r.Start(), "want range start time %q, got %q")
	assertSameString(t, "12:00", r.End(), "want range end time %q, got %q")
}

func TestAddMoreRanges(t *testing.T) {
	s := scheduler.NewSchedule()
	s, _ = s.AddRange("09:00", "12:00")
	s, _ = s.AddRange("15:00", "18:00")

	assertSameLength(t, 2, len(s))

	r := s[1]
	assertSameString(t, "15:00", r.Start(), "want range start time %q, got %q")
	assertSameString(t, "18:00", r.End(), "want range end time %q, got %q")
}

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name       string
		times      [][2]string
		wantLength int
		wantTimes  [][2]string
	}{
		{
			"in order, separate times",
			[][2]string{
				{"09:00", "12:00"},
				{"18:00", "20:00"},
			},
			2,
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
			1,
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
			1,
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
			1,
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
			1,
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
			2,
			[][2]string{
				{"09:00", "12:00"},
				{"15:00", "19:00"},
			},
		},
		{
			"separate time, added not in order, overlapping",
			[][2]string{
				{"20:00", "22:00"},
				{"15:00", "18:00"},
				{"09:00", "12:00"},
				{"16:00", "19:00"},
				{"14:00", "18:00"},
			},
			3,
			[][2]string{
				{"09:00", "12:00"},
				{"14:00", "19:00"},
				{"20:00", "22:00"},
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

			assertSameLength(t, tt.wantLength, len(s))

			for i, ts := range tt.wantTimes {
				r := s[i]
				assertSameString(t, ts[0], r.Start(), "want range start %q, got %q")
				assertSameString(t, ts[1], r.End(), "want range end %q, got %q")
			}
		})
	}
}

func TestRangeStartLessThanEnd(t *testing.T) {
	_, err := scheduler.NewSchedule().AddRange("14:00", "09:00")
	if err == nil {
		t.Errorf("expect error, because given start time is greater than end time, got %v", err)
	}

	_, err = scheduler.NewSchedule().AddRange("14:00", "14:00")
	if err == nil {
		t.Errorf("expect error, because given start time is equal to end time, got %v", err)
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
