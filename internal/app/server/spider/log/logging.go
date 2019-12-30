package log

import (
	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
)

func Debug(label string, action string, context box.Dict) {
	log.Debug(jsonic.MustEncode(box.Dict{
		"module":  "spider",
		"label":   label,
		"action":  action,
		"context": context,
	}))
}

func Info(label string, action string, context box.Dict) {
	log.Info(jsonic.MustEncode(box.Dict{
		"module":  "spider",
		"label":   label,
		"action":  action,
		"context": context,
	}))
}

func Result(value interface{}, context box.Dict) {
	log.Info(jsonic.MustEncode(box.Dict{
		"module":  "spider",
		"label":   "result",
		"result":  value,
		"context": context,
	}))
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}
