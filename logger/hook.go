package logger

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ErrHookNotFound = errors.New("hook not found")

type (
	HookGenerator    struct{}
	HookGenerateFunc func(*hierarchy.Hierarchy) (logrus.Hook, error)
)

var hookGeneratorMap = map[string]HookGenerateFunc{}

func init() {
	i := &HookGenerator{}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		hookGeneratorMap[strings.ToLower(method.Name)] = func(h *hierarchy.Hierarchy) (logrus.Hook, error) {
			in := []reflect.Value{v, reflect.ValueOf(h)}
			out := method.Func.Call(in)

			if !out[1].IsNil() {
				return nil, out[1].Interface().(error)
			}

			return out[0].Interface().(logrus.Hook), nil
		}
	}
}

func NewHook(typ string, c *hierarchy.Hierarchy) (logrus.Hook, error) {
	if hookGenerator, ok := hookGeneratorMap[typ]; ok {
		return hookGenerator(c)
	}

	return nil, fmt.Errorf("%w: %s", ErrHookNotFound, typ)
}

var _ logrus.Hook = &FileHook{}

// FileHook stores the hook of rolling file appender.
type FileHook struct {
	// formatter logrus.Formatter
	logger        *lumberjack.Logger
	logLevels     *LogLevels
	formatOptions *FormatOptions
}

// Fire is called when a log event is fired.
func (h *FileHook) Fire(entry *logrus.Entry) error {
	var (
		line []byte
		err  error
	)

	if entry.Context == nil {
		entry.Context = context.WithValue(context.Background(), "formatOptions", h.formatOptions)
	} else {
		entry.Context = context.WithValue(entry.Context, "formatOptions", h.formatOptions)
	}

	line, err = entry.Bytes()
	if err != nil {
		return err
	}

	// Write the the logger
	_, err = h.logger.Write(line)
	if err != nil {
		return err
	}

	return nil
}

// Levels returns the available logging
func (h *FileHook) Levels() []logrus.Level {
	return h.logLevels.ToLogrus()
}

func (*HookGenerator) File(h *hierarchy.Hierarchy) (logrus.Hook, error) {
	logger := &lumberjack.Logger{
		Filename:   h.GetString("file"),   // {var} is replaced
		MaxSize:    h.GetInt("size"),      // megabytes
		MaxBackups: h.GetInt("backup"),    // backup count
		MaxAge:     h.GetInt("days"),      // days
		Compress:   h.GetBool("compress"), // disabled by default
	}

	return &FileHook{logger, NewLogLevels(h), NewFormatOptions(h)}, nil
}

var stderr = colorable.NewColorableStderr()

var _ logrus.Hook = &StderrHook{}

// StderrHook is for stdout.
type StderrHook struct {
	writer.Hook
	logLevels     *LogLevels
	formatOptions *FormatOptions
}

func (h *StderrHook) Fire(entry *logrus.Entry) error {
	var (
		line []byte
		err  error
	)

	if entry.Context == nil {
		entry.Context = context.WithValue(context.Background(), "formatOptions", h.formatOptions)
	} else {
		entry.Context = context.WithValue(entry.Context, "formatOptions", h.formatOptions)
	}

	line, err = entry.Bytes()
	if err != nil {
		return err
	}

	_, err = h.Writer.Write(line)
	return err
}

// Levels returns the available logging
func (h *StderrHook) Levels() []logrus.Level {
	return h.logLevels.ToLogrus()
}

func (*HookGenerator) Stderr(h *hierarchy.Hierarchy) (logrus.Hook, error) {
	hook := writer.Hook{
		Writer:    stderr,
		LogLevels: logrus.AllLevels,
	}

	logLevels := NewLogLevels(h)
	formatOptions := NewFormatOptions(h)

	return &StderrHook{hook, logLevels, formatOptions}, nil
}

var stdout = colorable.NewColorableStdout()

var _ logrus.Hook = &StdoutHook{}

// StdoutHook is for stdout.
type StdoutHook struct {
	writer.Hook
	logLevels     *LogLevels
	formatOptions *FormatOptions
}

// Levels returns the available logging
func (h *StdoutHook) Levels() []logrus.Level {
	return h.logLevels.ToLogrus()
}

func (h *StdoutHook) Fire(entry *logrus.Entry) error {
	var (
		line []byte
		err  error
	)

	if entry.Context == nil {
		entry.Context = context.WithValue(context.Background(), "formatOptions", h.formatOptions)
	} else {
		entry.Context = context.WithValue(entry.Context, "formatOptions", h.formatOptions)
	}

	line, err = entry.Bytes()
	if err != nil {
		return err
	}

	_, err = h.Writer.Write(line)
	return err
}

func (*HookGenerator) Stdout(h *hierarchy.Hierarchy) (logrus.Hook, error) {
	hook := writer.Hook{
		Writer:    stdout,
		LogLevels: logrus.AllLevels,
	}

	return &StdoutHook{hook, NewLogLevels(h), NewFormatOptions(h)}, nil
}

var _ logrus.Hook = &TelegramHook{}

type TelegramHook struct {
	router        *router.ServiceRouter
	logLevels     *LogLevels
	formatOptions *FormatOptions
}

const telegramURL = "telegram://%s@telegram?chats=%s"

func (h *TelegramHook) Fire(entry *logrus.Entry) error {
	var (
		line []byte
		err  error
	)

	if entry.Context == nil {
		entry.Context = context.WithValue(context.Background(), "formatOptions", h.formatOptions)
	} else {
		entry.Context = context.WithValue(entry.Context, "formatOptions", h.formatOptions)
	}

	line, err = entry.Bytes()
	if err != nil {
		return err
	}

	errs := h.router.Send(string(line), nil)
	if len(errs) > 0 && errs[0] != nil {
		return errs[0]
	}

	return nil
}

func (h *TelegramHook) Levels() []logrus.Level {
	return h.logLevels.ToLogrus()
}

func (*HookGenerator) Telegram(h *hierarchy.Hierarchy) (logrus.Hook, error) {
	router, err := shoutrrr.CreateSender(fmt.Sprintf(telegramURL, h.GetString("token"), h.GetString("chat_id")))
	if err != nil {
		return nil, err
	}

	return &TelegramHook{router, NewLogLevels(h), NewFormatOptions(h)}, nil
}
