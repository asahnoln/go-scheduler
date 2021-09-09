package scheduler

import (
	"bufio"
	"io"
	"regexp"
	"time"
)

type Items map[string]Week

// CLI holds scanner to scan user stdin
type CLI struct {
	scanner *bufio.Scanner
}

// NewCLI creates a new CLI for given reader
func NewCLI(r io.Reader) CLI {
	return CLI{
		bufio.NewScanner(r),
	}
}

// Process check given reader, scans the command and create a Week
func (c CLI) Process() (Items, error) {
	ws := make(Items)
	re := regexp.MustCompile(`^add\s+(\w+)\s+(\w+)\s+(\d{2}:\d{2})-(\d{2}:\d{2})$`)

	// TODO: Check for !ok
	for c.scanner.Scan() {
		ms := re.FindSubmatch([]byte(c.scanner.Text()))

		s, err := NewSchedule().Add(string(ms[3]), string(ms[4]))
		if err != nil {
			return ws, err
		}

		person := string(ms[1])
		w := ws[person]
		d := parseDay(string(ms[2]))
		ws[person] = w.Add(d, s)
	}

	return ws, nil
}

func (i Items) Item(name string) Week {
	return i[name]
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
