package signals

import (
	"os"
	"strings"
)

func Join(sig []os.Signal, sep string) string {
	var names []string
	for _, s := range sig {
		names = append(names, s.String())
	}
	return strings.Join(names, sep)
}
