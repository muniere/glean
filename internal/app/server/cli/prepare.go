package cli

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/muniere/glean/internal/pkg/lumber"
)

func prepare(options *options) error {
	if options.verbose {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if len(options.logDir) == 0 {
		if err := prepareForConsoleLog(options); err != nil {
			return err
		}
		base := &logrus.TextFormatter{
			DisableColors:    false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "15:04:05.000",
		}
		logrus.SetFormatter(&lumber.TextFormatter{base})
	} else {
		if err := prepareForFileLog(options); err != nil {
			return err
		}
		base := &logrus.JSONFormatter{
			TimestampFormat:  "15:04:05.000",
			DisableTimestamp: false,
			DataKey:          "fields",
			FieldMap:         nil,
			CallerPrettyfier: nil,
			PrettyPrint:      false,
		}
		logrus.SetFormatter(&lumber.JSONFormatter{base})
	}

	return nil
}

func prepareForConsoleLog(options *options) error {
	logrus.SetOutput(os.Stderr)
	return nil
}

func prepareForFileLog(options *options) error {
	var err error

	logrus.SetOutput(ioutil.Discard)

	err = prepareCmdFileLog(options)
	if err != nil {
		return err
	}

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

func prepareCmdFileLog(options *options) error {
	file, err := os.OpenFile(
		path.Join(options.logDir, cmdLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHookWithFilter(file, logrus.AllLevels, func(entry *logrus.Entry) bool {
		return entry.Data["command"] != nil
	}))
	return nil
}

func prepareOutFileLog(options *options) error {
	file, err := os.OpenFile(
		path.Join(options.logDir, outLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHook(file, logrus.AllLevels))
	return nil
}

func prepareErrFileLog(options *options) error {
	file, err := os.OpenFile(
		path.Join(options.logDir, errLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHook(file, []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}))
	return nil
}
