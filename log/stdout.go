package log

import (
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var stdout = colorable.NewColorableStdout()

// StdoutHook is for stdout.
type StdoutHook struct {
	writer.Hook
	logLevels *LogLevels
}

// Levels returns the available logging levels
func (sh *StdoutHook) Levels() []logrus.Level {
	return sh.LogLevels
}

func (*HookGenerator) Stdout(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	logLevels := NewLogLevels()
	if err := logLevels.ReadConfig(c); err != nil {
		return nil, err
	}

	hook := writer.Hook{
		Writer:    stdout,
		LogLevels: logrus.AllLevels,
	}

	return &StdoutHook{hook, logLevels}, nil
}
