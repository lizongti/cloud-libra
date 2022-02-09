package hook

import (
	"io"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/sirupsen/logrus"
)

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
