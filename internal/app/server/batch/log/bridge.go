package log

import (
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
)

var (
	Warn  = logrus.Warn
	Error = logrus.Error
)

func Start(context box.Dict) {
	logrus.WithFields(t(context, "", "start")).Info()
}

func StartStep(step string, context box.Dict) {
	logrus.WithFields(t(context, step, "start")).Info()
}

func Finish(context box.Dict) {
	logrus.WithFields(t(context, "", "finish")).Info()
}

func FinishStep(step string, context box.Dict) {
	logrus.WithFields(t(context, step, "finish")).Info()
}

func Skip(context box.Dict) {
	logrus.WithFields(t(context, "", "skip")).Info()
}

func SkipStep(step string, context box.Dict) {
	logrus.WithFields(t(context, step, "skip")).Info()
}

func Result(value interface{}, context box.Dict) {
	logrus.WithFields(t(context, "", "result")).WithField("result", value).Info()
}

func ResultStep(step string, value interface{}, context box.Dict) {
	logrus.WithFields(t(context, step, "result")).WithField("result", value).Info()
}

func t(context box.Dict, step string, suffix string) logrus.Fields {
	pc, file, line, _ := runtime.Caller(2)

	x := logrus.Fields{
		"module": "batch",
		"file": path.Join(path.Base(path.Dir(file)), path.Base(file)),
		"line": line,
		"context": context,
	}

	if len(step) > 0 {
		x["event"] = strings.Join([]string{step, suffix}, "::")
	} else {
		fun := runtime.FuncForPC(pc)
		elems := strings.Split(fun.Name(), ".")
		step := elems[len(elems)-1]
		x["event"] = strings.Join([]string{step, suffix}, "::")
	}

	return x
}
