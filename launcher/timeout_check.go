// +build go1.15

package launcher

import (
	"errors"
	"os"
)

func isDeadlineExceededErr(err error) bool {
	return errors.Is(err, os.ErrDeadlineExceeded)
}
