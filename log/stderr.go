package log

import (
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var stderr = colorable.NewColorableStderr()

// StderrHook is for stdout.
type StderrHook struct {
	writer.Hook
	logLevels *LogLevels
}

// Levels returns the available logging levels
func (sh *StderrHook) Levels() []logrus.Level {
	return sh.LogLevels
}

func (*HookGenerator) Stderr(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	logLevels := NewLogLevels()
	if err := logLevels.ReadConfig(c); err != nil {
		return nil, err
	}

	hook := writer.Hook{
		Writer:    stderr,
		LogLevels: logrus.AllLevels,
	}
	return &StderrHook{hook, logLevels}, nil
}
