package lumber

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/pathname"
	"github.com/muniere/glean/internal/pkg/std"
)

func Warn(err error) {
	logrus.WithFields(w(std.NewDict(std.Pair("error", err.Error())))).Warn()
}

func Error(err error) {
	logrus.WithFields(w(std.NewDict(std.Pair("error", err.Error())))).Error()
}

func Start(ctx std.Dict) {
	logrus.WithFields(t("", "start", ctx)).Info()
}

func StartStep(step string, ctx std.Dict) {
	logrus.WithFields(t(step, "start", ctx)).Info()
}

func Finish(ctx std.Dict) {
	logrus.WithFields(t("", "finish", ctx)).Info()
}

func FinishStep(step string, ctx std.Dict) {
	logrus.WithFields(t(step, "finish", ctx)).Info()
}

func Skip(ctx std.Dict) {
	logrus.WithFields(t("", "skip", ctx)).Info()
}

func SkipStep(step string, ctx std.Dict) {
	logrus.WithFields(t(step, "skip", ctx)).Info()
}

func Result(v interface{}, ctx std.Dict) {
	logrus.WithFields(t("", "result", ctx)).WithField("result", v).Info()
}

func ResultStep(step string, v interface{}, ctx std.Dict) {
	logrus.WithFields(t(step, "result", ctx)).WithField("result", v).Info()
}

func t(step string, suffix string, ctx std.Dict) logrus.Fields {
	pc, file, line, _ := runtime.Caller(2)

	x := logrus.Fields{
		"module":  "batch",
		"file":    pathname.MustLeaf(file, 2),
		"line":    line,
		"context": ctx.Values(),
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

func w(ctx std.Dict) logrus.Fields {
	_, file, line, _ := runtime.Caller(2)

	return logrus.Fields{
		"module":  "batch",
		"file":    pathname.MustLeaf(file, 2),
		"line":    line,
		"context": ctx.Values(),
	}
}
