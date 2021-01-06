// +build !go1.15

package launcher

import (
	"os"
)

func isDeadlineExceededErr(err error) bool {
	return err != nil && os.IsTimeout(err)
}
