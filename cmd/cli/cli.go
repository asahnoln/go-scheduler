package cli

import (
	"io"
	"regexp"
	"time"

	"github.com/asahnoln/go-scheduler"
)

type CLI struct {
	r io.Reader
}

func NewCLI(r io.Reader) CLI {
	return CLI{r}
}

func (c CLI) Process() (scheduler.Week, error) {
	w := scheduler.NewWeek()

	cmd, err := io.ReadAll(c.r)
	if err != nil {
		return w, err
	}

	re := regexp.MustCompile(`^\s*add\s+(\w+)\s+(\d{2}:\d{2})-(\d{2}:\d{2})\s*$`)
	ms := re.FindSubmatch([]byte(cmd))

	s := scheduler.NewSchedule()
	s, err = s.Add(string(ms[2]), string(ms[3]))

	return w.Add(time.Monday, s), err
}
