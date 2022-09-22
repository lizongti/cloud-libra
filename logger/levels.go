package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnknownLogLevel = errors.New("unknown log level")
	ErrLevelNotSet     = errors.New("log level not set")
)

type LogLevel struct {
	level logrus.Level
}

func NewLogLevel(i interface{}) *LogLevel {
	switch i := i.(type) {
	case string:
		s := i

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
			panic(fmt.Errorf("%w: %s", ErrUnknownLogLevel, s))
		}

		return &LogLevel{level}
	case *hierarchy.Hierarchy:
		h := i

		return NewLogLevel(h.GetString("level"))

	default:
		panic(fmt.Errorf("%w: %T", ErrUnknownLogLevel, i))
	}
}

func (l *LogLevel) ToLogrus() logrus.Level {
	return l.level
}

type LogLevels struct {
	logLevels []*LogLevel
}

func NewLogLevels(i interface{}) *LogLevels {
	switch i := i.(type) {
	case string:
		str := i

		logLevel := NewLogLevel(str)
		logLevels := make([]*LogLevel, 0, logLevel.ToLogrus()+1)

		for i := logrus.Level(0); i <= logLevel.ToLogrus(); i++ {
			s := logrus.Level(i).String()
			logLevels = append(logLevels, NewLogLevel(s))
		}

		return &LogLevels{logLevels}

	case []string:
		strs := i

		logLevels := make([]*LogLevel, 0, len(strs))
		for _, str := range strs {
			logLevels = append(logLevels, NewLogLevel(str))
		}

		return &LogLevels{logLevels}

	case *hierarchy.Hierarchy:
		h := i

		var a any
		if h.IsArray("level") {
			a = h.GetStringSlice("level")

			return NewLogLevels(a)
		}

		a = h.GetString("level")

		return NewLogLevels(a)

	default:
		panic(fmt.Errorf("%w: %T", ErrUnknownLogLevel, i))
	}
}

func (l *LogLevels) ToLogrus() []logrus.Level {
	logrusLevels := make([]logrus.Level, 0, len(l.logLevels))
	for _, level := range l.logLevels {
		logrusLevels = append(logrusLevels, level.ToLogrus())
	}

	return logrusLevels
}
