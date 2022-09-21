package formatter

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var _ logrus.Formatter = (*TextFormatter)(nil)

type TextFormatter struct {
	logrus.TextFormatter
}

func (*Generator) Text(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
	formatter := &TextFormatter{
		TextFormatter: logrus.TextFormatter{
			ForceColors:            true,
			TimestampFormat:        "2006/01/02 15:04:05.0000000", // the "time" field configuratiom
			FullTimestamp:          true,
			DisableLevelTruncation: true, // log upgrade field configuration
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := f.File
				fileIndex := strings.LastIndex(s, "/")
				packageIndex := strings.LastIndex(s[:fileIndex], "/")
				index := strings.LastIndex(s[packageIndex+1:fileIndex], "@")
				var packageFile string

				if index >= 0 {
					packageFile = s[packageIndex+1:]
				} else {
					packageFile = s[packageIndex+1 : fileIndex][index+1:] + s[fileIndex:]
				}

				file := fmt.Sprintf(" %s:%d:", packageFile, f.Line)
				function := fmt.Sprintf("[%s]", f.Function[strings.LastIndex(f.Function, ".")+1:])

				return function, file
			},
		},
	}

	return formatter, nil
}
