package log

import (
	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/lumber"
)

func Debug(action string, context box.Dict) {
	lumber.Debug(box.Dict{
		"module":  "batch",
		"action":  action,
		"context": context,
	})
}

func Info(action string, context box.Dict) {
	lumber.Info(box.Dict{
		"module":  "batch",
		"action":  action,
		"context": context,
	})
}

func Result(value interface{}, context box.Dict) {
	lumber.Info(box.Dict{
		"module":  "batch",
		"result":  value,
		"context": context,
	})
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}
