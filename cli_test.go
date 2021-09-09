package scheduler_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler"
)

func TestCLIAddSchedulesSeparately(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:00-12:00\nadd apollo monday 14:00-15:00")

	c := scheduler.NewCLI(in, nil)
	ws, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	s := ws.Item("apollo").Day(time.Monday)

	assertSameLength(t, 2, len(s))
	assertSameString(t, "09:00", s[0].StartString(), "want start time %q, got %q")
	assertSameString(t, "14:00", s[1].StartString(), "want start time %q, got %q")
}

func TestCLIAddDays(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:30-12:00\nadd apollo tuesday 14:30-15:00")

	c := scheduler.NewCLI(in, nil)
	ws, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	w := ws.Item("apollo")
	s1 := w.Day(time.Monday)
	s2 := w.Day(time.Tuesday)

	assertSameLength(t, 1, len(s1))
	assertSameLength(t, 1, len(s2))
	assertSameString(t, "09:30", s1[0].StartString(), "want start time %q, got %q")
	assertSameString(t, "14:30", s2[0].StartString(), "want start time %q, got %q")
}

func TestCLIAddItems(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:00-12:00\nadd arthur monday 14:00-15:00")

	c := scheduler.NewCLI(in, nil)
	ws, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	t.Run("get arthur schedule", func(t *testing.T) {
		w := ws.Item("arthur")
		s := w.Day(time.Monday)

		assertSameLength(t, 1, len(s))
		assertSameString(t, "14:00", s[0].StartString(), "want start time %q, got %q")
	})

	t.Run("get apollo schedule", func(t *testing.T) {
		w := ws.Item("apollo")
		s := w.Day(time.Monday)

		assertSameLength(t, 1, len(s))
		assertSameString(t, "09:00", s[0].StartString(), "want start time %q, got %q")
	})
}

func TestCliWholeWeek(t *testing.T) {
	days := map[time.Weekday]string{
		time.Monday:    "monday",
		time.Tuesday:   "tuesday",
		time.Wednesday: "wednesday",
		time.Thursday:  "thursday",
		time.Friday:    "friday",
		time.Saturday:  "saturday",
		time.Sunday:    "sunday",
	}

	for i, d := range days {
		t.Run(d, func(t *testing.T) {
			in := strings.NewReader(fmt.Sprintf("add apollo %s 09:15-12:00", d))
			c := scheduler.NewCLI(in, nil)
			ws, err := c.Process()
			assertNoError(t, err, "unexpected error while processing CLI: %v")

			s := ws.Item("apollo").Day(i)
			assertSameLength(t, 1, len(s))
			assertSameString(t, "09:15", s[0].StartString(), "want start time %q, got %q")

		})
	}
}

func TestCLIWrongCommand(t *testing.T) {
	in := strings.NewReader("this is wrong")

	c := scheduler.NewCLI(in, nil)

	_, err := c.Process()
	assertError(t, err, "want CLI error, because wrong command, got: %v")
}

func TestCLIShowData(t *testing.T) {
	in := strings.NewReader("add apollo monday 09:00-14:00\nadd apollo monday 15:00-18:00\nshow apollo")
	out := bytes.Buffer{}

	c := scheduler.NewCLI(in, &out)

	_, err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	want := `
apollo

Monday
09:00-14:00
15:00-18:00
`[1:]

	assertSameString(t, want, out.String(), "want output\n%v\n\ngot\n%v")
}
