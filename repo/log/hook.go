package log

import (
	"errors"
	"io"
	"strings"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ErrUnsupportedLogHookCreater = errors.New("unsupported log hook creater")
	ErrUnsupportedLogLevelSource = errors.New("unsupported log level source")
	ErrUnsupportedLogLevel       = errors.New("unsupported log level")
)

type HookCreater interface {
	Create(*tree.MapTree) (logrus.Hook, error)
	String() string
}

var hookCreaterMap = map[string]HookCreater{}

func RegisterHookCreater(hookCreater HookCreater) {
	hookCreaterMap[hookCreater.String()] = hookCreater
}

func NewHook(name string, config *tree.MapTree) (logrus.Hook, error) {
	name = strings.ToLower(name)
	hookCreater, ok := hookCreaterMap[name]
	if !ok {
		return nil, ErrUnsupportedLogHookCreater
	}
	return hookCreater.Create(config)
}

type LumberjackHook struct {
	logger   *lumberjack.Logger
	logLevel *LogLevel
}

func (hook *LumberjackHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	_, err = hook.logger.Write(line)
	return err
}

func (hook *LumberjackHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type LumberjackHookCreater struct{}

func (c LumberjackHookCreater) Create(config *tree.MapTree) (logrus.Hook, error) {
	hook := &LumberjackHook{}
	// var directory string
	var file = cast.ToString(config.Get(magic.UnixChain("file")))
	var filename string
	switch cast.ToString(config.Get(magic.UnixChain("directory"))) {
	case "":
		filename = file
	case "exe":
	case "cwd":
	case "project":
	}
	hook.logger = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cast.ToInt(config.Get(magic.UnixChain("size"))),      // megabytes
		MaxBackups: cast.ToInt(config.Get(magic.UnixChain("backup"))),    // backups
		MaxAge:     cast.ToInt(config.Get(magic.UnixChain("day"))),       // days
		Compress:   cast.ToBool(config.Get(magic.UnixChain("compress"))), // disabled by default
	}

	logLevel, err := NewLogLevel(c.levelSource(config))
	if err != nil {
		return nil, err
	}
	hook.logLevel = logLevel

	return hook, nil
}

func (LumberjackHookCreater) String() string {
	return "lumberjack"
}

func (LumberjackHookCreater) levelSource(config *tree.MapTree) interface{} {
	level := cast.ToString(config.Get(magic.UnixChain("level")))
	if level != "" {
		return level
	}
	levels := cast.ToStringSlice(config.Get(magic.UnixChain("levels")))
	if len(levels) > 0 {
		return levels
	}
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

type StdoutHook struct {
	writer   io.Writer
	logLevel *LogLevel
}

func (hook *StdoutHook) Fire(entry *logrus.Entry) error {
	line, out := entry.Bytes()
	if out != nil {
		return out
	}
	_, out = hook.writer.Write(line)
	return out
}

func (hook *StdoutHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type StdoutHookCreater struct{}

func (c StdoutHookCreater) Create(config *tree.MapTree) (logrus.Hook, error) {
	hook := &StdoutHook{}
	hook.writer = stdio.Out()

	logLevel, out := NewLogLevel(c.levelSource(config))
	if out != nil {
		return nil, out
	}
	hook.logLevel = logLevel

	return hook, nil
}

func (StdoutHookCreater) String() string {
	return "stdout"
}

func (StdoutHookCreater) levelSource(config *tree.MapTree) interface{} {
	level := cast.ToString(config.Get(magic.UnixChain("level")))
	if level != "" {
		return level
	}
	levels := cast.ToStringSlice(config.Get(magic.UnixChain("levels")))
	if len(levels) > 0 {
		return levels
	}
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func init() {
	RegisterHookCreater(StdoutHookCreater{})
}

type StderrHook struct {
	writer   io.Writer
	logLevel *LogLevel
}

func (hook *StderrHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(line)
	return err
}

func (hook *StderrHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type StderrHookCreater struct{}

func (c StderrHookCreater) Create(config *tree.MapTree) (logrus.Hook, error) {
	hook := &StderrHook{}
	hook.writer = stdio.Err()

	logLevel, err := NewLogLevel(c.levelSource(config))
	if err != nil {
		return nil, err
	}
	hook.logLevel = logLevel

	return hook, nil
}

func (StderrHookCreater) String() string {
	return "stderr"
}

func (StderrHookCreater) levelSource(config *tree.MapTree) interface{} {
	level := cast.ToString(config.Get(magic.UnixChain("level")))
	if level != "" {
		return level
	}
	levels := cast.ToStringSlice(config.Get(magic.UnixChain("levels")))
	if len(levels) > 0 {
		return levels
	}
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

func init() {
	RegisterHookCreater(StderrHookCreater{})
}
