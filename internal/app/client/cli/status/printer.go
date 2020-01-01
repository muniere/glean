package status

import (
	"fmt"
	"strings"
)

func printLine(values ...string) {
	fmt.Println(strings.Join(values, "\t"))
}
