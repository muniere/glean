package sys

import (
	"os"
)

func CheckError(err error) {
	if err == nil {
		return
	}

	os.Exit(1)
}
