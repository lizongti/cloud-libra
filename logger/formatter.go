package logger

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var (
	ErrFormatterNotFound = errors.New("formatter not found")
	ErrMethodNotValid    = errors.New("method not valid")
)

type (
	FormatterGenerator    struct{}
	FormatterGenerateFunc func(*hierarchy.Hierarchy) (logrus.Formatter, error)
)

var formatterGeneratorMap = map[string]FormatterGenerateFunc{}

func init() {
	i := &FormatterGenerator{}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		formatterGeneratorMap[strings.ToLower(method.Name)] = func(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
			in := []reflect.Value{v, reflect.ValueOf(h)}
			out := method.Func.Call(in)

			if !out[1].IsNil() {
				return nil, out[1].Interface().(error)
			}

			return out[0].Interface().(logrus.Formatter), nil
		}
	}
}

func NewFormatter(c *hierarchy.Hierarchy) (logrus.Formatter, error) {
	typ := c.GetString("type")
	if fn, ok := formatterGeneratorMap[typ]; ok {
		return fn(c)
	}

	return nil, fmt.Errorf("%w: %s", ErrFormatterNotFound, typ)
}

var ErrUnexpectedFieldKey = errors.New("unexpected field key")

var _ logrus.Formatter = (*JSONFormatter)(nil)

type JSONFormatter struct {
	logrus.JSONFormatter
}

func (tf *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	bytes, ok := entry.Data["Bytes"]
	if ok {
		return bytes.([]byte), nil
	}

	bytes, err := tf.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	entry.Data["Bytes"] = bytes
	return bytes.([]byte), nil
}

func (*FormatterGenerator) JSON(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
	fieldMap := make(logrus.FieldMap)

	for k, v := range h.GetStringMapString("fieldMap") {
		switch k {
		case logrus.FieldKeyMsg:
			fieldMap[logrus.FieldKeyMsg] = v
		case logrus.FieldKeyLevel:
			fieldMap[logrus.FieldKeyLevel] = v
		case logrus.FieldKeyTime:
			fieldMap[logrus.FieldKeyTime] = v
		case logrus.FieldKeyLogrusError:
			fieldMap[logrus.FieldKeyLogrusError] = v
		case logrus.FieldKeyFunc:
			fieldMap[logrus.FieldKeyFunc] = v
		case logrus.FieldKeyFile:
			fieldMap[logrus.FieldKeyFile] = v
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnexpectedFieldKey, k)
		}
	}

	formatter := &JSONFormatter{
		JSONFormatter: logrus.JSONFormatter{
			DisableTimestamp:  true,
			DisableHTMLEscape: true,
			FieldMap:          fieldMap,
			PrettyPrint:       false,
		},
	}

	return formatter, nil
}

var _ logrus.Formatter = (*TextFormatter)(nil)

type TextFormatter struct {
	logrus.TextFormatter
}

func (tf *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	bytes, ok := entry.Data["Bytes"]
	if ok {
		return bytes.([]byte), nil
	}

	bytes, err := tf.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	entry.Data["Bytes"] = bytes
	return bytes.([]byte), nil
}

func (*FormatterGenerator) Text(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
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
				atIndex := strings.LastIndex(s[packageIndex+1:fileIndex], "@")
				var packageFile string
				if atIndex >= 0 {
					packageFile = s[packageIndex+1:]
				} else {
					packageFile = s[packageIndex+1 : fileIndex][atIndex+1:] + s[fileIndex:]
				}

				funcIndex := strings.LastIndex(f.Function, ".")
				structIndex := strings.LastIndex(f.Function[:funcIndex], ".")
				var function string
				if structIndex >= 0 {
					function = f.Function[structIndex+1:]
				} else {
					function = f.Function[funcIndex+1:]
				}

				return fmt.Sprintf("%s:", function), fmt.Sprintf(" %s:%d", packageFile, f.Line)
			},
		},
	}

	return formatter, nil
}
