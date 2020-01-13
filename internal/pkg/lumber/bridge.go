package lumber

import (
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/pathname"
	"github.com/muniere/glean/internal/pkg/std"
)

func Trace(values std.Dict) {
	logrus.WithFields(t(values)).Trace()
}

func Debug(values std.Dict) {
	logrus.WithFields(t(values)).Debug()
}

func Info(values std.Dict) {
	logrus.WithFields(t(values)).Info()
}

func Warn(values std.Dict) {
	logrus.WithFields(t(values)).Warn()
}

func Error(values std.Dict) {
	logrus.WithFields(t(values)).Error()
}

func Panic(values std.Dict) {
	logrus.WithFields(t(values)).Panic()
}

func Fatal(values std.Dict) {
	logrus.WithFields(t(values)).Fatal()
}

func t(values std.Dict) logrus.Fields {
	x := logrus.Fields{}
	for k, v := range values.Values() {
		x[k] = v
	}

	vf := values.Has("file")
	vl := values.Has("line")

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
