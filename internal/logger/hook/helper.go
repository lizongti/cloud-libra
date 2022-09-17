package hook

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func parseLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info", "print":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		panic(fmt.Errorf("unknown type of log level %s", level))
	}
}

func aboveLevel(level string) []logrus.Level {
	var logrusLevels []logrus.Level
	for i := logrus.Level(0); i <= parseLevel(level); i++ {
		logrusLevels = append(logrusLevels, i)
	}
	return logrusLevels
}

func parseLevels(levels []string) []logrus.Level {
	var logrusLevels = make([]logrus.Level, 0, len(levels))
	for _, level := range levels {
		logrusLevels = append(logrusLevels, parseLevel(level))
	}
	return logrusLevels
}
