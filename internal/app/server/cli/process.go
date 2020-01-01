package cli

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/lumber"
)

func prepare(options *options) error {
	if options.verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if len(options.logDir) == 0 {
		if err := prepareForConsoleLog(options); err != nil {
			return err
		}
		base := &log.TextFormatter{
			DisableColors:    false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "15:04:05.000",
		}
		log.SetFormatter(&lumber.TextFormatter{base})
	} else {
		if err := prepareForFileLog(options); err != nil {
			return err
		}
		base := &log.JSONFormatter{
			TimestampFormat:  "15:04:05.000",
			DisableTimestamp: false,
			DataKey:          "values",
			FieldMap:         nil,
			CallerPrettyfier: nil,
			PrettyPrint:      false,
		}
		log.SetFormatter(&lumber.JSONFormatter{base})
	}

	return nil
}

func prepareForConsoleLog(options *options) error {
	log.SetOutput(os.Stderr)
	return nil
}

func prepareForFileLog(options *options) error {
	var err error

	log.SetOutput(ioutil.Discard)

	err = prepareOutFileLog(options)
	if err != nil {
		return err
	}

	err = prepareErrFileLog(options)
	if err != nil {
		return err
	}

	return nil
}

func prepareOutFileLog(options *options) error {
	file, err := os.OpenFile(
		path.Join(options.logDir, outLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	log.AddHook(lumber.NewFileHook(file, log.AllLevels))
	return nil
}

func prepareErrFileLog(options *options) error {
	file, err := os.OpenFile(
		path.Join(options.logDir, errLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	log.AddHook(lumber.NewFileHook(file, []log.Level{
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	}))
	return nil
}
