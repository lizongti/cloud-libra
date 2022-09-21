package logger

import (
	"io"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

func New(h *hierarchy.Hierarchy) (*logrus.Logger, error) {
	logger := logrus.New()

	// Disable output.
	logger.SetOutput(io.Discard)

	// Enable or disable caller.
	if h.GetBool("caller") {
		logger.SetReportCaller(true)
	}

	// Set formatter.
	formatter, err := NewFormatter(h.Sub("formatter"))
	if err != nil {
		return nil, err
	}
	logger.SetFormatter(formatter)

	// Set hooks.
	h.ForeachInArray("hooks", func(index int, hierarchy *hierarchy.Hierarchy) (bool, error) {
		typ := hierarchy.GetString("type")
		hook, err := NewHook(typ, hierarchy)
		if err != nil {
			return false, err
		}
		logger.AddHook(hook)
		return true, nil
	})

	// Set Level
	logLevel, err := NewLogLevel(h.GetString("level"))
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel.ToLogrus())

	return logger, nil
}
