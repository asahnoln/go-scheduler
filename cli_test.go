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
quit
`[1:])
	out := bytes.Buffer{}

	c := scheduler.NewCLI(in, &out)

	_ = c.MainLoop()
	// assertNoError(t, err, "unexpected error while processing CLI: %v")

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

func TestCLIAddDifferentItems(t *testing.T) {
	names := []string{
		"arthur", "apollo",
	}

	for _, n := range names {
		t.Run("get "+n+" schedule", func(t *testing.T) {
			in := strings.NewReader("add " + n + " thursday 14:25-15:15\nshow " + n + "\nquit")
			out := &bytes.Buffer{}

			c := scheduler.NewCLI(in, out)
			_ = c.MainLoop()
			// assertNoError(t, err, "unexpected error while processing CLI: %v")

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
			in := strings.NewReader(fmt.Sprintf("add apollo %s 09:15-12:00\nshow apollo\nquit", d))
			out := bytes.Buffer{}
			c := scheduler.NewCLI(in, &out)
			_ = c.MainLoop()
			// assertNoError(t, err, "unexpected error while processing CLI: %v")

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
	in := strings.NewReader("this is wrong\nquit")
	out := &bytes.Buffer{}

	c := scheduler.NewCLI(in, out)

	_ = c.MainLoop()
	assertSameString(t, "unknown command string \"this is wrong\"\n", out.String(), "want CLI error %q, got %q")
}

func TestCLIShowSeveralPeople(t *testing.T) {
	in := strings.NewReader(`
add apollo monday 09:00-14:00
add arthur monday 10:00-15:00
add apollo tuesday 12:00-13:00
add arthur friday 18:45-19:45
show apollo arthur
quit
`[1:])
	out := bytes.Buffer{}

	c := scheduler.NewCLI(in, &out)

	_ = c.MainLoop()
	// assertNoError(t, err, "unexpected error while processing CLI: %v")

	want := `
apollo arthur

Monday
09:00-15:00

Tuesday
12:00-13:00

Friday
18:45-19:45

`[1:]

	assertSameString(t, want, out.String(), "want output\n%v\n\ngot\n%v")
}

func TestCLIQuitOnlyOnCommand(t *testing.T) {
	cmds := []string{"show apollo", "add apollo monday 19:00-21:00", "wrong"}

	for _, cmd := range cmds {
		t.Run(cmd+" should run", func(t *testing.T) {
			in := strings.NewReader(cmd)
			out := &bytes.Buffer{}
			c := scheduler.NewCLI(in, out)
			quit := make(chan error)

			go func() {
				quit <- c.MainLoop()
			}()

			select {
			case <-time.After(1 * time.Millisecond):
			case err := <-quit:
				t.Errorf("want continue loop, got loop break on command %q with error: %v", cmd, err)
			}
		})
	}
}

func TestCLIQuitOnCommand(t *testing.T) {
	in := strings.NewReader("quit")
	out := &bytes.Buffer{}
	c := scheduler.NewCLI(in, out)
	quit := make(chan error)

	go func() {
		quit <- c.MainLoop()
	}()

	select {
	case <-time.After(1 * time.Millisecond):
		t.Errorf("want quit loop, got loop continuation on \"quit\" command")
	case <-quit:
	}
}
