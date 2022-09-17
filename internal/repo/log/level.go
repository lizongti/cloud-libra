package log

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogLevel struct {
	levels []logrus.Level
}

func NewLogLevel(v interface{}) (l *LogLevel, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			l = nil
		}
	}()
	l = &LogLevel{}
	switch v := v.(type) {
	case string:
		level, err := l.parse(v)
		if err != nil {
			return nil, err
		}
		for i := logrus.Level(0); i <= level; i++ {
			l.levels = append(l.levels, i)
		}
	case []string:
		for _, s := range v {
			level, err := l.parse(s)
			if err != nil {
				return nil, err
			}
			l.levels = append(l.levels, level)
		}
	case logrus.Level:
		for i := logrus.Level(0); i <= v; i++ {
			l.levels = append(l.levels, i)
		}
	case []logrus.Level:
		l.levels = make([]logrus.Level, len(v))
		copy(l.levels, v)
	default:
		return nil, ErrUnsupportedLogLevelSource
	}

	sort.Slice(l.levels, func(i, j int) bool {
		return l.levels[i] < l.levels[j]
	})

	return l, nil
}

func (l LogLevel) Level() logrus.Level {
	return l.levels[len(l.levels)-1]
}

func (l LogLevel) Levels() []logrus.Level {
	return l.levels
}

func (LogLevel) parse(level string) (logrus.Level, error) {
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
		return 0, ErrUnsupportedLogLevel
	}
}
