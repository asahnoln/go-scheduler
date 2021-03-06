package scheduler_test

import (
	"testing"

	"github.com/asahnoln/go-scheduler"
	"github.com/asahnoln/go-scheduler/test"
)

func TestNewRange(t *testing.T) {
	got, err := scheduler.NewRangeFromStrings("09:00", "12:00")
	test.AssertNoError(t, err, "unexpected error while creating range: %v")

	test.AssertSameString(t, "09:00", got.StartString(), "want range start time %q, got %q")
	test.AssertSameString(t, "12:00", got.EndString(), "want range end time %q, got %q")
}

func TestCreateSchedule(t *testing.T) {
	want, err := scheduler.NewRangeFromStrings("09:00", "12:00")
	test.AssertNoError(t, err, "unexpected error while creating range: %v")

	s := scheduler.NewSchedule().AddRange(want)
	test.AssertSameLength(t, 1, len(s))

	got := s[0]
	test.AssertSameRange(t, want, got)
}

func TestAddMoreRanges(t *testing.T) {
	s := scheduler.NewSchedule()
	s, _ = s.Add("09:00", "12:00")
	s, _ = s.Add("15:00", "18:00")

	test.AssertSameLength(t, 2, len(s))

	r := s[1]
	test.AssertSameString(t, "15:00", r.StartString(), "want range start time %q, got %q")
	test.AssertSameString(t, "18:00", r.EndString(), "want range end time %q, got %q")
}

// TODO: Assert Error objects
func TestRangeStartLessThanEnd(t *testing.T) {
	_, err := scheduler.NewSchedule().Add("14:00", "09:00")
	test.AssertError(t, err, "expect error, because given start time is greater than end time, got %v")

	_, err = scheduler.NewSchedule().Add("14:00", "14:00")
	test.AssertError(t, err, "expect error, because given start time is equal to end time, got %v")
}

// TODO: Test for Error objects
func TestWrongTimes(t *testing.T) {
	_, err := scheduler.NewSchedule().Add("wrong", "15:00")
	test.AssertError(t, err, "want error because start time is wrong, got %v")

	_, err = scheduler.NewSchedule().Add("15:00", "wrong")
	test.AssertError(t, err, "want error because end time is wrong, got %v")
}
