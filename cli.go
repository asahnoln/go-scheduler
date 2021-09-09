package scheduler

import (
	"bufio"
	"io"
	"regexp"
	"time"
)

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
func (c CLI) Process() (Week, error) {
	w := NewWeek()
	re := regexp.MustCompile(`^add\s+\w+\s+(\w+)\s+(\d{2}:\d{2})-(\d{2}:\d{2})$`)

	// TODO: Check for !ok
	for c.scanner.Scan() {
		ms := re.FindSubmatch([]byte(c.scanner.Text()))

		s, err := NewSchedule().Add(string(ms[2]), string(ms[3]))
		if err != nil {
			return w, err
		}

		var d time.Weekday
		switch string(ms[1]) {
		case "monday":
			d = time.Monday
		case "tuesday":
			d = time.Tuesday
		}

		w = w.Add(d, s)
	}

	return w, nil
}
