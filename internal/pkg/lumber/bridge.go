package lumber

import (
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/pathname"
)

func Trace(values box.Dict) {
	logrus.WithFields(t(values)).Trace()
}

func Debug(values box.Dict) {
	logrus.WithFields(t(values)).Debug()
}

func Info(values box.Dict) {
	logrus.WithFields(t(values)).Info()
}

func Warn(values box.Dict) {
	logrus.WithFields(t(values)).Warn()
}

func Error(values box.Dict) {
	logrus.WithFields(t(values)).Error()
}

func Panic(values box.Dict) {
	logrus.WithFields(t(values)).Panic()
}

func Fatal(values box.Dict) {
	logrus.WithFields(t(values)).Fatal()
}

func t(values box.Dict) logrus.Fields {
	x := logrus.Fields{}
	for k, v := range values {
		x[k] = v
	}

	_, vf := values["file"]
	_, vl := values["line"]

	if vf && vl {
		return x
	}

	_, file, line, _ := runtime.Caller(2)

	if !vf {
		x["file"] = pathname.MustLeaf(file, 2)
	}
	if !vl {
		x["line"] = line
	}

	return x
}
