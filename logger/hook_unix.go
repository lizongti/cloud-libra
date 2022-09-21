//go:build !windows && !nacl && !plan9
// +build !windows,!nacl,!plan9

package logger

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	*syslog.SyslogHook
	logLevels *levels.LogLevels
}

// Levels returns the available logging levels
func (sh *SyslogHook) Levels() []logrus.Level {
	return sh.logLevels.ToLogrus()
}

func (*HookFactory) Syslog(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	var a any
	if c.IsArray("level") {
		a = c.GetStringSlice("level")
	} else {
		a = c.GetString("level")
	}

	logLevels, err := levels.NewLogLevels(a)
	if err != nil {
		return nil, err
	}

	syslogHook := syslog.NewSyslogHook(
		c.GetString("network"),
		c.GetString("addr"),
		syslog.Log_DEBUG,
		c.GetString("tag"),
	),

	hook := &SysLogHook{
		SyslogHook: syslogHook,
		logLevels:  logLevels,
	}

	return &SyslogHook{hook, logLevels}, nil
}