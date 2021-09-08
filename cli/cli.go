package cli

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/asahnoln/go-scheduler"
)

// CLI holds scanner to scan user stdin
type CLI struct {
	in *bufio.Scanner
}

// NewCLI creates a new CLI for given reader
func NewCLI(r io.Reader) CLI {
	return CLI{
		bufio.NewScanner(r),
	}
}

// Process check given reader, scans the command and create a Week
func (c CLI) Process() (scheduler.Week, error) {
	w := scheduler.NewWeek()

	// TODO: Check for !ok
	_ = c.in.Scan()

	re := regexp.MustCompile(`^\s*add\s+(\w+)\s+(\d{2}:\d{2})-(\d{2}:\d{2})\s*$`)
	ms := re.FindSubmatch([]byte(c.in.Text()))

	s := scheduler.NewSchedule()
	s, err := s.Add(string(ms[2]), string(ms[3]))

	return w.Add(time.Monday, s), err
}
