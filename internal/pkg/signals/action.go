package signals

import (
	"os"
	"os/signal"
)

func Wait(sig ...os.Signal) os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	return <-ch
}
