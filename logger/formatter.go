package logger

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var (
	ErrFormatterNotFound     = errors.New("formatter not found")
	ErrMethodNotValid        = errors.New("method not valid")
	ErrFormatOptionsNotFound = errors.New("format options not found")
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

func NewFormatter(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
	typ := h.GetStringVal("type", "text")
	if fn, ok := formatterGeneratorMap[typ]; ok {
		return fn(h)
	}

	return nil, fmt.Errorf("%w: %s", ErrFormatterNotFound, typ)
}

var _ logrus.Formatter = (*TextFormatter)(nil)

type TextFormatter struct {
	logrus.TextFormatter
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry.Context == nil || entry.Context.Value("formatOptions") == nil {
		return f.TextFormatter.Format(entry)
	}

	c, ok := entry.Context.Value("formatOptions").(*FormatOptions)
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrFormatOptionsNotFound, entry.Context.Value("formatOptions"))
	}

	if !c.date && !c.time && !c.nanosecond && !c.timezone {
		f.TextFormatter.DisableTimestamp = true
	} else {
		f.TextFormatter.DisableTimestamp = false
		f.TextFormatter.TimestampFormat = f.generateTimeFormat(c.date, c.time, c.nanosecond, c.timezone)
	}

	f.TextFormatter.CallerPrettyfier = f.generateCallerPrettierfier(c.file, c.function)

	data, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	if !c.level {
		index := bytes.Index(data, []byte("[0m"))
		data = data[index+4:]
	}

	return data, nil
}

func (f *TextFormatter) generateTimeFormat(date bool, time bool, nanosecond bool, timezone bool) string {
	var buf bytes.Buffer
	if date {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}

		buf.WriteString("2006/01/02")
	}

	if time {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}

		buf.WriteString("15:04:05")
	}

	if nanosecond {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}

		buf.WriteString(".000000000")
	}

	if timezone {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}

		buf.WriteString("-0700")
	}

	return buf.String()
}

func (f *TextFormatter) generateCallerPrettierfier(file bool, function bool) func(*runtime.Frame) (string, string) {
	switch {
	case file && function:
		return func(frame *runtime.Frame) (string, string) {
			return f.generateFile(frame) + f.generateFunction(frame) + ":", ""
		}
	case file:
		return func(frame *runtime.Frame) (string, string) {
			return f.generateFile(frame) + ":", ""
		}
	case function:
		return func(frame *runtime.Frame) (string, string) {
			return f.generateFunction(frame) + ":", ""
		}
	default:
		return func(frame *runtime.Frame) (string, string) {
			return "", ""
		}
	}
}

func (f *TextFormatter) generateFile(frame *runtime.Frame) string {
	s := frame.File
	fileIndex := strings.LastIndex(s, "/")
	packageIndex := strings.LastIndex(s[:fileIndex], "/")
	atIndex := strings.LastIndex(s[packageIndex+1:fileIndex], "@")

	if atIndex >= 0 {
		return fmt.Sprintf(" %s:%d", s[packageIndex+1:], frame.Line)
	}

	return fmt.Sprintf(" %s%s:%d", s[packageIndex+1 : fileIndex][atIndex+1:], s[fileIndex:], frame.Line)
}

func (f *TextFormatter) generateFunction(frame *runtime.Frame) string {
	s := frame.Function

	return fmt.Sprintf(" (%s)", s[strings.LastIndex(s, "/")+1:])
}

func (*FormatterGenerator) Text(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
	formatter := &TextFormatter{
		TextFormatter: logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
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

var ErrUnexpectedFieldKey = errors.New("unexpected field key")

var _ logrus.Formatter = (*JSONFormatter)(nil)

type JSONFormatter struct {
	logrus.JSONFormatter
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	bytes, ok := entry.Data["Bytes"]
	if ok {
		return bytes.([]byte), nil
	}

	bytes, err := f.JSONFormatter.Format(entry)
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

var ErrUnexpectedFormat = errors.New("unexpected format")

type FormatOptions struct {
	level      bool
	date       bool
	time       bool
	nanosecond bool
	timezone   bool
	file       bool
	function   bool
	message    bool
}

func NewFormatOptions(h *hierarchy.Hierarchy) *FormatOptions {
	c := &FormatOptions{}

	if h.IsArray("format") {
		for _, option := range h.GetStringSlice("format") {
			switch option {
			case "level", "lvl", "lv":
				c.level = true
			case "date":
				c.date = true
			case "time", "tm":
				c.time = true
			case "nano-second", "nanosecond", "nano", "ns":
				c.nanosecond = true
			case "timezone", "tz":
				c.timezone = true
			case "file", "f":
				c.file = true
			case "function", "func", "fn":
				c.function = true
			case "message", "msg", "m":
				c.message = true
			}
		}
	} else {
		switch format := h.GetString("format"); format {
		case "all":
			c.level = true
			c.date = true
			c.time = true
			c.nanosecond = true
			c.timezone = true
			c.file = true
			c.function = true
			c.message = true
		case "default":
			c.level = true
			c.date = true
			c.time = true
			c.file = true
			c.message = true
		case "simple":
			c.level = true
			c.date = true
			c.time = true
			c.message = true
		case "message", "":
			c.message = true
		default:
			panic(fmt.Errorf("%w: %s", ErrUnexpectedFormat, format))
		}
	}

	return c
}
