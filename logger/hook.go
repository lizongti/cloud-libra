package logger

import (
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
	logger    *lumberjack.Logger
	logLevels *LogLevels
}

// Fire is called when a log event is fired.
func (lh *FileHook) Fire(entry *logrus.Entry) error {
	// Convert the line to string
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	// Write the the logger
	_, err = lh.logger.Write(line)
	if err != nil {
		return err
	}

	return nil
}

// Levels returns the available logging
func (lh *FileHook) Levels() []logrus.Level {
	return lh.logLevels.ToLogrus()
}

func (*HookGenerator) File(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	var a any
	if c.IsArray("level") {
		a = c.GetStringSlice("level")
	} else {
		a = c.GetString("level")
	}

	logLevels, err := NewLogLevels(a)
	if err != nil {
		return nil, err
	}

	logger := &lumberjack.Logger{
		Filename:   c.GetString("file"),   // {var} is replaced
		MaxSize:    c.GetInt("size"),      // megabytes
		MaxBackups: c.GetInt("backup"),    // backup count
		MaxAge:     c.GetInt("days"),      // days
		Compress:   c.GetBool("compress"), // disabled by default
	}

	return &FileHook{logger, logLevels}, nil
}

var stderr = colorable.NewColorableStderr()

var _ logrus.Hook = &StderrHook{}

// StderrHook is for stdout.
type StderrHook struct {
	writer.Hook
	logLevels *LogLevels
}

// Levels returns the available logging
func (sh *StderrHook) Levels() []logrus.Level {
	return sh.logLevels.ToLogrus()
}

func (*HookGenerator) Stderr(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	var a any
	if c.IsArray("level") {
		a = c.GetStringSlice("level")
	} else {
		a = c.GetString("level")
	}

	logLevels, err := NewLogLevels(a)
	if err != nil {
		return nil, err
	}

	hook := writer.Hook{
		Writer:    stderr,
		LogLevels: logrus.AllLevels,
	}

	return &StderrHook{hook, logLevels}, nil
}

var stdout = colorable.NewColorableStdout()

var _ logrus.Hook = &StdoutHook{}

// StdoutHook is for stdout.
type StdoutHook struct {
	writer.Hook
	logLevels *LogLevels
}

// Levels returns the available logging
func (sh *StdoutHook) Levels() []logrus.Level {
	return sh.logLevels.ToLogrus()
}

func (*HookGenerator) Stdout(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	var a any
	if c.IsArray("level") {
		a = c.GetStringSlice("level")
	} else {
		a = c.GetString("level")
	}

	logLevels, err := NewLogLevels(a)
	if err != nil {
		return nil, err
	}

	hook := writer.Hook{
		Writer:    stdout,
		LogLevels: logrus.AllLevels,
	}

	return &StdoutHook{hook, logLevels}, nil
}

var _ logrus.Hook = &TelegramHook{}

type TelegramHook struct {
	router    *router.ServiceRouter
	logLevels *LogLevels
}

const telegramURL = "telegram://%s@telegram?chats=%s"

func (th *TelegramHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	errs := th.router.Send(string(line), nil)
	if len(errs) > 0 && errs[0] != nil {
		return errs[0]
	}

	return nil
}

func (th *TelegramHook) Levels() []logrus.Level {
	return th.logLevels.ToLogrus()
}

func (*HookGenerator) Telegram(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	var a any
	if c.IsArray("level") {
		a = c.GetStringSlice("level")
	} else {
		a = c.GetString("level")
	}

	logLevels, err := NewLogLevels(a)
	if err != nil {
		return nil, err
	}

	router, err := shoutrrr.CreateSender(fmt.Sprintf(telegramURL, c.GetString("token"), c.GetString("chat_id")))
	if err != nil {
		return nil, err
	}

	return &TelegramHook{router, logLevels}, nil
}
