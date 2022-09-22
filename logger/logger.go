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

	// Enable caller.
	logger.SetReportCaller(true)

	// Set formatter.
	formatter, err := NewFormatter(h.Sub("formatter"))
	if err != nil {
		return nil, err
	}
	logger.SetFormatter(formatter)

	// Set hooks.
	h.ForeachInArray("hooks", func(index int, h *hierarchy.Hierarchy) (bool, error) {
		typ := h.GetString("type")
		hook, err := NewHook(typ, h)
		if err != nil {
			return false, err
		}
		logger.AddHook(hook)
		return true, nil
	})

	// Set Level
	logger.SetLevel(NewLogLevel(h.GetStringVal("level", "info")).ToLogrus())

	return logger, nil
}
