package lumber

import (
	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
)

func Trace(values box.Dict) {
	log.WithFields(t(values)).Trace()
}

func Debug(values box.Dict) {
	log.WithFields(t(values)).Debug()
}

func Info(values box.Dict) {
	log.WithFields(t(values)).Info()
}

func Warn(values box.Dict) {
	log.WithFields(t(values)).Warn()
}

func Error(values box.Dict) {
	log.WithFields(t(values)).Error()
}

func Panic(values box.Dict) {
	log.WithFields(t(values)).Panic()
}

func Fatal(values box.Dict) {
	log.WithFields(t(values)).Fatal()
}

func t(values box.Dict) log.Fields {
	x := log.Fields{}
	for k, v := range values {
		x[k] = v
	}
	return x
}
