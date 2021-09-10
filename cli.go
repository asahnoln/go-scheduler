package scheduler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

type items map[string]Week

// CLI holds scanner to scan user stdin
type CLI struct {
	scanner *bufio.Scanner
	out     io.Writer
}

var days = [7]time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

// NewCLI creates a new CLI for given reader
func NewCLI(r io.Reader, w io.Writer) CLI {
	return CLI{
		bufio.NewScanner(r),
		w,
	}
}

// MainLoop check given reader, scans the command and create a Week
func (c CLI) MainLoop() error {
	ws := make(items)

	for {
		for c.scanner.Scan() {
			words := strings.Split(c.scanner.Text(), " ")
			if len(words) < 1 {
				return fmt.Errorf("unknown command verb, got %q", c.scanner.Text())
			}

			switch words[0] {
			case "show":
				c.show(ws)
			case "add":
				args := words[1:]
				if len(args) < 3 {
					return fmt.Errorf("not enough params, got %q", c.scanner.Text())
				}

				s := NewSchedule()
				var err error
				for _, timeRange := range args[2:] {
					r := strings.Split(timeRange, "-")
					s, err = s.Add(r[0], r[1])
					if err != nil {
						return err
					}
				}

				person := args[0]
				w := ws[person]
				d := parseDay(args[1])
				ws[person] = w.Add(d, s)
			case "quit":
				return fmt.Errorf("Quitted")
			default:
				fmt.Fprintf(c.out, "unknown command string %q\n", c.scanner.Text())
			}

		}
	}
}

func (c CLI) show(ws items) {
	args := strings.Split(c.scanner.Text(), " ")[1:]
	fmt.Fprintf(c.out, "%s\n\n", strings.Join(args, " "))

	w := NewWeek()
	for _, d := range days {
		s := NewSchedule()
		for _, a := range args {
			s = s.AddSchedule(ws[a].Day(d))
		}
		w = w.Add(d, s)
	}

	for _, d := range days {
		s := w.Day(d)
		if len(s) == 0 {
			continue
		}

		fmt.Fprintf(c.out, "%s\n", time.Weekday(d))
		for _, r := range s {
			fmt.Fprintf(c.out, "%v-%v\n", r.StartString(), r.EndString())
		}
		fmt.Fprintf(c.out, "\n")
	}
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
