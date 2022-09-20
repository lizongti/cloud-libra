package log

import (
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LumberjackHook stores the hook of rolling file appender.
type LumberjackHook struct {
	logger    *lumberjack.Logger
	logLevels *LogLevels
}

// Fire is called when a log event is fired.
func (lh *LumberjackHook) Fire(entry *logrus.Entry) error {
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

// Levels returns the available logging levels.
func (lh *LumberjackHook) Levels() []logrus.Level {
	return lh.logLevels.Levels()
}

func (*HookGenerator) Lumberjack(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	logLevels := NewLogLevels()
	if err := logLevels.ReadConfig(c); err != nil {
		return nil, err
	}

	logger := &lumberjack.Logger{
		Filename:   c.GetString("file"),   // {var} is replaced
		MaxSize:    c.GetInt("size"),      // megabytes
		MaxBackups: c.GetInt("backup"),    // backup count
		MaxAge:     c.GetInt("days"),      // days
		Compress:   c.GetBool("compress"), // disabled by default
	}

	return &LumberjackHook{logger, logLevels}, nil
}
