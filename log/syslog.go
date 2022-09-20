//go:build !windows && !nacl && !plan9
// +build !windows,!nacl,!plan9

package log

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
)

// SyslogConfig stores the configuration of SyslogHook
type SyslogConfig struct {
	Network string
	Addr    string
	Tag     string
	Level   string
	Levels  []string
}

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	*syslog.SyslogHook
	logLevels *LogLevels
}

// Levels returns the available logging levels
func (sh *SyslogHook) Levels() []logrus.Level {
	return sh.LogLevels
}

func (*HookFactory) Syslog(c *hierarchy.Hierarchy) (logrus.Hook, error) {
	logLevels := NewLogLevels()
	if err := logLevels.ReadConfig(c); err != nil {
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