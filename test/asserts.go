package test

import (
	"testing"

	"github.com/asahnoln/go-scheduler"
)

func AssertError(t testing.TB, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatalf(message, err)
	}
}

func AssertNoError(t testing.TB, err error, message string) {
	t.Helper()

	if err != nil {
		t.Fatalf(message, err)
	}
}

func AssertSameLength(t testing.TB, want, got int) {
	t.Helper()

	if want != got {
		t.Fatalf("want range count %d, got %d", want, got)
	}
}

func AssertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}

func AssertSameRange(t testing.TB, want, got scheduler.Range) {
	t.Helper()

	if want != got {
		t.Errorf("want range %v-%v, got %v-%v", want.StartString(), want.EndString(), got.StartString(), got.EndString())
	}
}
