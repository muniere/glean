package lumber

import (
	"os"

	"github.com/sirupsen/logrus"
)

type FileHook struct {
	file   *os.File
	levels []logrus.Level
}

func NewFileHook(file *os.File, levels []logrus.Level) FileHook {
	return FileHook{
		file:   file,
		levels: levels,
	}
}

func (hook FileHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook FileHook) Fire(entry *logrus.Entry) error {
	if entry.Logger.Out == os.Stdout || entry.Logger.Out == os.Stderr {
		return nil
	}

	msg, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.file.Write(msg)
	return err
}
