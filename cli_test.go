package scheduler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler"
)

func TestCLIAddSchedulesSeparately(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:00-12:00\nadd apollo monday 14:00-15:00")

	c := scheduler.NewCLI(in)
	w, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	s := w.Day(time.Monday)

	assertSameLength(t, 2, len(s))
	assertSameString(t, "09:00", s[0].StartString(), "want start time %q, got %q")
	assertSameString(t, "14:00", s[1].StartString(), "want start time %q, got %q")
}

func TestCLIAddDays(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:30-12:00\nadd apollo tuesday 14:30-15:00")

	c := scheduler.NewCLI(in)
	w, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	s1 := w.Day(time.Monday)
	s2 := w.Day(time.Tuesday)

	assertSameLength(t, 1, len(s1))
	assertSameLength(t, 1, len(s2))
	assertSameString(t, "09:30", s1[0].StartString(), "want start time %q, got %q")
	assertSameString(t, "14:30", s2[0].StartString(), "want start time %q, got %q")
}
