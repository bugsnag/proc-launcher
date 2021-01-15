// +build !go1.15

package launcher

import (
	"testing"
	"fmt"
)

type timeout interface {
	Timeout() bool
}

type someErr struct{ isTimeout bool }

func (e someErr) Error() string {
	return "bad"
}

func (e someErr) Timeout() bool {
	return e.isTimeout
}

func TestIsDeadlineExceeded(t *testing.T) {
	err1 := someErr{ true }
	if !isDeadlineExceededErr(err1) {
		t.Errorf("Failed to detect timeout error")
	}
	err2 := someErr{ false }
	if isDeadlineExceededErr(err2) {
		t.Errorf("Detected non-timeout error as a timeout")
	}
	err3 := fmt.Errorf("some other error")
	if isDeadlineExceededErr(err3) {
		t.Errorf("Detected non-timeout error as a timeout")
	}
}
