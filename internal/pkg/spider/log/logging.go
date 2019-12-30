package log

import (
	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/box"
	"github.com/muniere/glean/internal/pkg/jsonic"
)

type dict = box.Dict

func Debug(label string, action string, context dict) {
	log.Debug(jsonic.MustEncode(dict{
		"label":   label,
		"action":  action,
		"context": context,
	}))
}

func Info(label string, action string, context dict) {
	log.Info(jsonic.MustEncode(dict{
		"label":   label,
		"action":  action,
		"context": context,
	}))
}

func Result(value interface{}, context dict) {
	log.Info(jsonic.MustEncode(dict{
		"label": "result",
		"result": value,
		"context": context,
	}))
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}
