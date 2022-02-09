package hook

import (
	"io"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/sirupsen/logrus"
)

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
