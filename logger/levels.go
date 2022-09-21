package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownLogLevel = errors.New("unknown log level")
	ErrLevelNotSet     = errors.New("log level not set")
)

type LogLevel struct {
	level logrus.Level
}

func NewLogLevel(s string) (*LogLevel, error) {
	var level logrus.Level
	switch strings.ToLower(s) {
	case "panic":
		level = logrus.PanicLevel
	case "fatal":
		level = logrus.FatalLevel
	case "error":
		level = logrus.ErrorLevel
	case "warn", "warning":
		level = logrus.WarnLevel
	case "info", "print":
		level = logrus.InfoLevel
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownLogLevel, level)
	}
	return &LogLevel{level}, nil
}

func (l *LogLevel) ToLogrus() logrus.Level {
	return l.level
}

type LogLevels struct {
	logLevels []*LogLevel
}

func NewLogLevels(i interface{}) (*LogLevels, error) {
	switch i := i.(type) {
	case string:
		str := i

		logLevel, err := NewLogLevel(str)
		if err != nil {
			return nil, err
		}

		logLevels := make([]*LogLevel, 0, logLevel.ToLogrus()+1)

		for i := logrus.Level(0); i <= logLevel.ToLogrus(); i++ {
			s := logrus.Level(i).String()

			logLevel, err := NewLogLevel(s)
			if err != nil {
				return nil, err
			}

			logLevels = append(logLevels, logLevel)
		}

		return &LogLevels{logLevels}, nil

	case []string:
		strs := i
		logLevels := make([]*LogLevel, 0, len(strs))

		for _, str := range strs {
			logLevel, err := NewLogLevel(str)
			if err != nil {
				return nil, err
			}

			logLevels = append(logLevels, logLevel)
		}

		return &LogLevels{logLevels}, nil

	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownLogLevel, i)
	}
}

func (l *LogLevels) ToLogrus() []logrus.Level {
	logrusLevels := make([]logrus.Level, 0, len(l.logLevels))
	for _, level := range l.logLevels {
		logrusLevels = append(logrusLevels, level.ToLogrus())
	}

	return logrusLevels
}
