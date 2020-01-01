package cli

import (
	"os"
	"os/signal"
	"strings"

	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
)

func wait(sig ...os.Signal) {
	lumber.Info(box.Dict{
		"module": "root",
		"action": "signal.wait",
		"values": join(sig, ", "),
	})

	defer lumber.Info(box.Dict{
		"module": "root",
		"action": "signal.recv",
		"values": join(sig, ", "),
	})

	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	<-ch
}

func join(sig []os.Signal, sep string) string {
	var names []string
	for _, s := range sig {
		names = append(names, s.String())
	}
	return strings.Join(names, sep)
}
