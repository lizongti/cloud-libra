package log

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ErrUnsupportedLogHookCreater = errors.New("unsupported log hook creature")
	ErrUnsupportedLogLevelSource = errors.New("unsupported log level source")
	ErrUnsupportedLogLevel       = errors.New("unsupported log level")
)

type Processor interface {
	Process([]byte) []byte
}

type Stdio struct {
	stdout io.Writer
	stderr io.Writer
	mutex  sync.Mutex
}

var stdio = Stdio{}

func (s *Stdio) Out() io.Writer {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stdout == nil {
		s.stdout = colorable.NewColorableStdout()
	}
	return s.stdout
}

func (s *Stdio) Err() io.Writer {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stderr == nil {
		s.stderr = colorable.NewColorableStderr()
	}
	return s.stderr
}

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
		for i := logrus.Level(0); i <= l.parse(v); i++ {
			l.levels = append(l.levels, i)
		}
	case []string:
		for _, level := range v {
			l.levels = append(l.levels, l.parse(level))
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
	return l, nil
}

func (l LogLevel) Levels() []logrus.Level {
	return l.levels
}

func (LogLevel) parse(level string) logrus.Level {
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
		panic(ErrUnsupportedLogLevel)
	}
}

func (LogLevel) aboveLevel(level logrus.Level) []logrus.Level {
	var levels []logrus.Level
	for i := logrus.Level(0); i <= level; i++ {
		levels = append(levels, i)
	}
	return levels
}

type HookCreater interface {
	Create(*tree.MapTree, Processor) (logrus.Hook, error)
	String() string
}

var hookCreaters = map[string]HookCreater{}

func register(hookCreater HookCreater) {
	hookCreaters[hookCreater.String()] = hookCreater
}

func NewHook(name string, config *tree.MapTree, processor Processor) (logrus.Hook, error) {
	name = strings.ToLower(name)
	hookCreater, ok := hookCreaters[name]
	if !ok {
		return nil, ErrUnsupportedLogHookCreater
	}
	return hookCreater.Create(config, processor)
}

type LumberjackHook struct {
	logger    *lumberjack.Logger
	processor Processor
	logLevel  *LogLevel
}

func (hook *LumberjackHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	if hook.processor != nil {
		line = hook.processor.Process(line)
	}

	_, err = hook.logger.Write(line)
	return err
}

func (hook *LumberjackHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type LumberjackHookCreater struct{}

func (c LumberjackHookCreater) Create(config *tree.MapTree, processor Processor) (logrus.Hook, error) {
	hook := &LumberjackHook{}
	hook.logger = &lumberjack.Logger{
		Filename:   cast.ToString(config.Get(magic.UnixChain("file"))),
		MaxSize:    cast.ToInt(config.Get(magic.UnixChain("size"))),      // megabytes
		MaxBackups: cast.ToInt(config.Get(magic.UnixChain("backup"))),    // backups
		MaxAge:     cast.ToInt(config.Get(magic.UnixChain("day"))),       // days
		Compress:   cast.ToBool(config.Get(magic.UnixChain("compress"))), // disabled by default
	}
	hook.processor = processor

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
	writer    io.Writer
	processor Processor
	logLevel  *LogLevel
}

func (hook *StdoutHook) Fire(entry *logrus.Entry) error {
	line, out := entry.Bytes()
	if out != nil {
		return out
	}
	if hook.processor != nil {
		line = hook.processor.Process(line)
	}
	_, out = hook.writer.Write(line)
	return out
}

func (hook *StdoutHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type StdoutHookCreater struct{}

func (c StdoutHookCreater) Create(config *tree.MapTree, processor Processor) (logrus.Hook, error) {
	hook := &StdoutHook{}
	hook.writer = stdio.Out()
	hook.processor = processor

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
	register(StdoutHookCreater{})
}

type StderrHook struct {
	writer    io.Writer
	processor Processor
	logLevel  *LogLevel
}

func (hook *StderrHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	if hook.processor != nil {
		line = hook.processor.Process(line)
	}
	_, err = hook.writer.Write(line)
	return err
}

func (hook *StderrHook) Levels() []logrus.Level {
	return hook.logLevel.Levels()
}

type StderrHookCreater struct{}

func (c StderrHookCreater) Create(config *tree.MapTree, processor Processor) (logrus.Hook, error) {
	hook := &StderrHook{}
	hook.writer = stdio.Err()
	hook.processor = processor

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
	register(StderrHookCreater{})
}
