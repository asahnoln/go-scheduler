package cli_test

import (
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-scheduler/cmd/cli"
)

func TestAddItem(t *testing.T) {
	in := strings.NewReader("add monday 09:00-12:00")

	cli := cli.NewCLI(in)

	w, err := cli.Process()
	if err != nil {
		t.Fatalf("unexpected error while processing CLI: %v", err)
	}

	s := w.Day(time.Monday)

	want := 1
	got := len(s)
	if want != got {
		t.Fatalf("want schedule count %d, got %d", want, got)
	}

	wantTime := "09:00"
	gotTime := s[0].StartString()
	if wantTime != gotTime {
		t.Errorf("want start time %q, got %q", want, got)
	}
}
