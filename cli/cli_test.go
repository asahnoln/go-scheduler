package cli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler/cli"
	"github.com/asahnoln/go-scheduler/test"
)

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
			c := cli.NewCLI(in, &out)
			err := c.MainLoop()
			test.AssertNoError(t, err, "unexpected error while processing CLI: %v")

			want := fmt.Sprintf(`
apollo

%s
09:15-12:00

`[1:], i)

			test.AssertSameString(t, want, out.String(), "want output\n%v\ngot\n%v")
		})
	}
}

func TestCLIWrongCommand(t *testing.T) {
	in := strings.NewReader("this is wrong\nquit")
	out := &bytes.Buffer{}

	c := cli.NewCLI(in, out)

	err := c.MainLoop()
	test.AssertNoError(t, err, "unexpected error while CLI main loop: %v")
	test.AssertSameString(t, "unknown command string \"this is wrong\"\n", out.String(), "want CLI error %q, got %q")
}

func TestCLIShowSeveralPeople(t *testing.T) {
	in := strings.NewReader(`
add apollo monday 09:00-14:00 16:00-18:00
add arthur monday 10:00-16:00
add apollo tuesday 12:00-13:00
add arthur friday 18:45-19:45
show apollo arthur
quit
`[1:])
	out := bytes.Buffer{}

	c := cli.NewCLI(in, &out)

	err := c.MainLoop()
	test.AssertNoError(t, err, "unexpected error while processing CLI: %v")

	want := `
apollo arthur

Monday
09:00-18:00

Tuesday
12:00-13:00

Friday
18:45-19:45

`[1:]

	test.AssertSameString(t, want, out.String(), "want output\n%vgot\n%v")
}

func TestCLIQuitOnlyOnCommand(t *testing.T) {
	cmds := []string{"show apollo", "add apollo monday 19:00-21:00", "wrong"}

	for _, cmd := range cmds {
		t.Run(cmd+" should run", func(t *testing.T) {
			in := strings.NewReader(cmd)
			out := &bytes.Buffer{}
			c := cli.NewCLI(in, out)
			finished := make(chan error)

			go func() {
				finished <- c.MainLoop()
			}()

			select {
			case <-time.After(1 * time.Millisecond):
			case err := <-finished:
				t.Errorf("want continue loop, got loop break on command %q with error: %v", cmd, err)
			}
		})
	}
}

func TestCLIQuitOnCommand(t *testing.T) {
	in := strings.NewReader("quit")
	out := &bytes.Buffer{}
	c := cli.NewCLI(in, out)
	finished := make(chan error)

	go func() {
		finished <- c.MainLoop()
	}()

	select {
	case <-time.After(1 * time.Millisecond):
		t.Errorf("want quit loop, got loop continuation on \"quit\" command")
	case <-finished:
	}
}

// func TestReadData(t *testing.T) {
// 	file, err := os.CreateTemp("", "db")
// 	test.AssertError(t, err, "error while creating temp file: %v")

// 	file.Write([]byte(`[]`)
// }

func TestSaveData(t *testing.T) {
	in := strings.NewReader(`
add apollo monday 09:00-14:00 16:00-18:00
add apollo tuesday 12:00-13:00
add arthur friday 18:45-19:45
save
quit
`[1:])
	out := &bytes.Buffer{}
	file, err := os.CreateTemp("", "db")
	test.AssertNoError(t, err, "error while creating temp file: %v")
	defer file.Close()
	defer os.Remove(file.Name())

	c := cli.NewCLI(in, out)
	c.DB(file)

	err = c.MainLoop()
	test.AssertNoError(t, err, "unexpected error while processing CLI: %v")

	file.Seek(0, 0)

	var got cli.Items
	json.NewDecoder(file).Decode(&got)

	want := cli.LastItems()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want saved items from CLI %v, got %v", want, got)
	}
}

func TestLoadData(t *testing.T) {
	// TODO: Refactor duplication
	in := strings.NewReader(`
add apollo monday 09:00-14:00 16:00-18:00
add apollo tuesday 12:00-13:00
add arthur friday 18:45-19:45
save
quit
`[1:])
	out := &bytes.Buffer{}
	file, _ := os.CreateTemp("", "db")
	defer file.Close()
	defer os.Remove(file.Name())

	c := cli.NewCLI(in, out)
	c.DB(file)

	_ = c.MainLoop()

	want := cli.LastItems()

	// Testing loading
	in = strings.NewReader(`
quit
`[1:])
	c = cli.NewCLI(in, out)
	c.DB(file)
	err := c.MainLoop()
	test.AssertNoError(t, err, "unexpected error while processing CLI: %v")

	got := cli.LastItems()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want saved items from CLI %v, got %v", want, got)
	}
}
