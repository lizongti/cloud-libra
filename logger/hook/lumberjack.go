package hook

import (
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/sirupsen/logrus"
)

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
