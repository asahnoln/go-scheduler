package scheduler_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler"
)

func TestCLIShowData(t *testing.T) {
	in := strings.NewReader(`
add apollo monday 09:00-14:00 15:00-18:00
add apollo monday 16:00-19:00
add apollo tuesday 12:00-13:00
show apollo
`[1:])
	out := bytes.Buffer{}

	c := scheduler.NewCLI(in, &out)

	err := c.Process()
	assertNoError(t, err, "unexpected error while processing CLI: %v")

	want := `
apollo

Monday
09:00-14:00
15:00-19:00

Tuesday
12:00-13:00

`[1:]

	assertSameString(t, want, out.String(), "want output\n%v\n\ngot\n%v")
}

func TestCLIAddItems(t *testing.T) {
	names := []string{
		"arthur", "apollo",
	}

	for _, n := range names {
		t.Run("get "+n+" schedule", func(t *testing.T) {
			in := strings.NewReader("add " + n + " thursday 14:25-15:15\nshow " + n)
			out := &bytes.Buffer{}

			c := scheduler.NewCLI(in, out)
			err := c.Process()
			assertNoError(t, err, "unexpected error while processing CLI: %v")

			want := n + `

Thursday
14:25-15:15

`
			assertSameString(t, want, out.String(), "want output\n%v\ngot\n%v")
		})
	}
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
			in := strings.NewReader(fmt.Sprintf("add apollo %s 09:15-12:00\nshow apollo", d))
			out := bytes.Buffer{}
			c := scheduler.NewCLI(in, &out)
			err := c.Process()
			assertNoError(t, err, "unexpected error while processing CLI: %v")

			want := fmt.Sprintf(`
apollo

%s
09:15-12:00

`[1:], i)

			assertSameString(t, want, out.String(), "want output\n%v\ngot\n%v")
		})
	}
}

func TestCLIWrongCommand(t *testing.T) {
	in := strings.NewReader("this is wrong")

	c := scheduler.NewCLI(in, nil)

	err := c.Process()
	assertError(t, err, "want CLI error, because wrong command, got: %v")
}
