package lumber

import (
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/pathname"
	. "github.com/muniere/glean/internal/pkg/stdlib"
)

func Trace(values Dict) {
	logrus.WithFields(t(values)).Trace()
}

func Debug(values Dict) {
	logrus.WithFields(t(values)).Debug()
}

func Info(values Dict) {
	logrus.WithFields(t(values)).Info()
}

func Warn(values Dict) {
	logrus.WithFields(t(values)).Warn()
}

func Error(values Dict) {
	logrus.WithFields(t(values)).Error()
}

func Panic(values Dict) {
	logrus.WithFields(t(values)).Panic()
}

func Fatal(values Dict) {
	logrus.WithFields(t(values)).Fatal()
}

func t(values Dict) logrus.Fields {
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
