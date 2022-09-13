package lowlevel

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetPackageFile gets package/file.go style return
func GetPackageFile(s string) string {
	fileIndex := strings.LastIndex(s, "/")
	packageIndex := strings.LastIndex(s[:fileIndex], "/")
	atIndex := strings.LastIndex(s[packageIndex+1:fileIndex], "@")
	if atIndex == -1 {
		return s[packageIndex+1:]
	}
	return s[packageIndex+1:packageIndex+atIndex+1] + "" + s[fileIndex:]
}

// LastPart splits s with sep, and get last piece
func LastPart(s string, sep string) string {
	lastIndex := strings.LastIndex(s, sep)
	if lastIndex < 0 {
		return s
	}
	return s[lastIndex+len(sep):]
}

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
