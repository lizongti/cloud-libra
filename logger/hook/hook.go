package hook

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/aceaura/libra/core/tree"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
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
