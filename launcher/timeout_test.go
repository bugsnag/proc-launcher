// +build go1.15

package launcher

import (
	"fmt"
	"os"
	"testing"
)

func TestDeadlineExceeded(t *testing.T) {
	err1 := os.ErrDeadlineExceeded
	if !isDeadlineExceededErr(err1) {
		t.Errorf("Failed to detect literal error")
	}
	err2 := fmt.Errorf("wrap timeout: %w", os.ErrDeadlineExceeded)
	if !isDeadlineExceededErr(err2) {
		t.Errorf("Failed to detect wrapped error")
	}
}
