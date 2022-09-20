package log

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownLevel = errors.New("unknown type of log level")
	ErrLevelNotSet  = errors.New("log level not set")
)

type LogLevels struct {
	levels []logrus.Level
}

func NewLogLevels() *LogLevels {
	return &LogLevels{}
}

func (l *LogLevels) Levels() []logrus.Level {
	return l.levels
}

func (l *LogLevels) parseLevel(level string) (logrus.Level, error) {
	switch strings.ToLower(level) {
	case "panic":
		return logrus.PanicLevel, nil
	case "fatal":
		return logrus.FatalLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	case "warn", "warning":
		return logrus.WarnLevel, nil
	case "info", "print":
		return logrus.InfoLevel, nil
	case "debug":
		return logrus.DebugLevel, nil
	case "trace":
		return logrus.TraceLevel, nil
	default:
		panic(fmt.Errorf("%w: %s", ErrUnknownLevel, level))
	}
}

func (l *LogLevels) ReadConfig(c *hierarchy.Hierarchy) error {
	if levelStr := c.GetString("level"); levelStr != "" {
		level, err := l.parseLevel(levelStr)
		if err != nil {
			return err
		}

		for i := logrus.Level(0); i <= level; i++ {
			l.levels = append(l.levels, i)
		}
	}
	if levelStrs := c.GetStringSlice("levels"); len(levelStrs) > 0 {
		for _, levelStr := range levelStrs {
			level, err := l.parseLevel(levelStr)
			if err != nil {
				return err
			}

			l.levels = append(l.levels, level)
		}
	}

	return ErrLevelNotSet
}
