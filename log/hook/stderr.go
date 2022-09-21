package hook

import (
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/cloudlibraries/libra/log/levels"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var stderr = colorable.NewColorableStderr()

// StderrHook is for stdout.
type StderrHook struct {
	writer.Hook
	logLevels *levels.LogLevels
}

// Levels returns the available logging levels.
func (sh *StderrHook) Levels() []logrus.Level {
	return sh.LogLevels
}

func (*Generator) Stderr(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	logLevels := levels.NewLogLevels()
	if err := logLevels.ReadConfig(c); err != nil {
		return nil, err
	}

	hook := writer.Hook{
		Writer:    stderr,
		LogLevels: logrus.AllLevels,
	}

	return &StderrHook{hook, logLevels}, nil
}
