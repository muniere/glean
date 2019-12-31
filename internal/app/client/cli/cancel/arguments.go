package cancel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func normalize(args []string) ([]int, error) {
	var ids []int
	var errs []string

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err == nil {
			ids = append(ids, id)
		} else {
			errs = append(errs, arg)
		}
	}

	if len(errs) > 0 {
		arg := strings.Join(errs, ", ")
		msg := fmt.Sprintf("values must be ID numbers: %v", arg)
		return nil, errors.New(msg)
	}

	return ids, nil
}
