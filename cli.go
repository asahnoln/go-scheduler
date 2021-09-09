package scheduler

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

type Items map[string]Week

// CLI holds scanner to scan user stdin
type CLI struct {
	scanner *bufio.Scanner
	out     io.Writer
}

// NewCLI creates a new CLI for given reader
func NewCLI(r io.Reader, w io.Writer) CLI {
	return CLI{
		bufio.NewScanner(r),
		w,
	}
}

// Process check given reader, scans the command and create a Week
func (c CLI) Process() (Items, error) {
	ws := make(Items)
	reCommand, err := regexp.Compile(`^(\w+)`)
	if err != nil {
		return ws, err
	}

	// TODO: Check for !ok
	for c.scanner.Scan() {
		ms := reCommand.FindSubmatch([]byte(c.scanner.Text()))
		if len(ms) < 1 {
			return ws, fmt.Errorf("unknown command verb, got %q", c.scanner.Text())
		}

		switch string(ms[1]) {
		case "show":
			return c.show(ws)
		case "add":
			re, err := regexp.Compile(`^add\s+(\w+)\s+(\w+)\s+(\d{2}:\d{2})-(\d{2}:\d{2})$`)
			ms := re.FindSubmatch([]byte(c.scanner.Text()))

			if len(ms) < 4 {
				return ws, fmt.Errorf("unknown command string, got %q", c.scanner.Text())
			}

			s, err := NewSchedule().Add(string(ms[3]), string(ms[4]))
			if err != nil {
				return ws, err
			}

			person := string(ms[1])
			w := ws[person]
			d := parseDay(string(ms[2]))
			ws[person] = w.Add(d, s)
		default:
			return ws, fmt.Errorf("unknown command string, got %q", c.scanner.Text())
		}

	}

	return ws, nil
}

// Item returns schedule for given name
func (i Items) Item(name string) Week {
	return i[name]
}

func (c CLI) show(ws Items) (Items, error) {
	re, err := regexp.Compile(`^show\s+(\w+)$`)
	if err != nil {
		return ws, err
	}

	ms := re.FindSubmatch([]byte(c.scanner.Text()))

	fmt.Fprintf(c.out, "%s\n\n", ms[1])

	for d, s := range ws.Item(string(ms[1])).days {
		if len(s) == 0 {
			continue
		}

		fmt.Fprintf(c.out, "%s\n", time.Weekday(d))

		for _, r := range s {
			fmt.Fprintf(c.out, "%v-%v\n", r.StartString(), r.EndString())
		}
	}

	return ws, nil
}

func parseDay(day string) time.Weekday {
	var d time.Weekday
	switch day {
	case "monday":
		d = time.Monday
	case "tuesday":
		d = time.Tuesday
	case "wednesday":
		d = time.Wednesday
	case "thursday":
		d = time.Thursday
	case "friday":
		d = time.Friday
	case "saturday":
		d = time.Saturday
	}

	return d
}
