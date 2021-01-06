// +build go1.15

package launcher

import (
	"os"
	"errors"
)

func isDeadlineExceededErr(err error) bool {
	return errors.Is(err, os.ErrDeadlineExceeded)
}
