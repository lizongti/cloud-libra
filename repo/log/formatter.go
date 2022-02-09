package log

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	ErrUnsupportedLogFormatterCreater = errors.New("unsupported log hook creater")
)

type FormatterCreater interface {
	Create() logrus.Formatter
	String() string
}

var formatterCreaterMap = map[string]FormatterCreater{}

func RegisterFormatterCreater(formatterCreater FormatterCreater) {
	formatterCreaterMap[formatterCreater.String()] = formatterCreater
}

func NewFormatterCreater(name string) (logrus.Formatter, error) {
	name = strings.ToLower(name)
	formatterCreater, ok := formatterCreaterMap[name]
	if !ok {
		return nil, ErrUnsupportedLogFormatterCreater
	}
	return formatterCreater.Create(), nil
}

type TextFormatterCreater struct{}

func (TextFormatterCreater) Create() logrus.Formatter {
	return &logrus.TextFormatter{
		ForceColors:            true,
		TimestampFormat:        "2006-01-02 15:04:05.999999999 -0700", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log upgrade field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileIndex := strings.LastIndex(f.File, "/")
			pkgIndex := strings.LastIndex(f.File[:fileIndex], "/")
			atIndex := strings.LastIndex(f.File[pkgIndex+1:fileIndex], "@")
			if atIndex == -1 {
				return "", fmt.Sprintf(" %s:%d", f.File[pkgIndex+1:], f.Line)
			}
			return "", fmt.Sprintf(" %s%s:%d", f.File[pkgIndex+1:pkgIndex+atIndex+1],
				f.File[fileIndex:], f.Line)
		},
	}
}

func (TextFormatterCreater) String() string {
	return "text"
}

func init() {
	RegisterFormatterCreater(TextFormatterCreater{})
}

type JSONFormatterCreater struct{}

func (JSONFormatterCreater) Create() logrus.Formatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05.999999999 -0700", // the "time" field configuratiom
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@msg",
			logrus.FieldKeyFunc:  "@func",
			logrus.FieldKeyFile:  "@file",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return f.Function, f.File
		},
		PrettyPrint: false,
	}
}

func (JSONFormatterCreater) String() string {
	return "json"
}

func init() {
	RegisterFormatterCreater(TextFormatterCreater{})
}
