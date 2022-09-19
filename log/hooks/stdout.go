package hooks

import (
	"encoding/json"
	"io"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var stdout = colorable.NewColorableStdout()

// StdoutConfig stores the configuration of StdoutHook.
type StdoutConfig struct {
	Level  string
	Levels []string
}

// StdoutHook is for stdout.
type StdoutHook struct {
	writer.Hook
	processor *Processor
}

type stdoutWriter struct {
	processor *Processor
	writer    io.Writer
}

func (s *stdoutWriter) Write(p []byte) (n int, err error) {
	if s.processor != nil && s.processor.Process != nil {
		p = s.processor.Process(p)
	}
	return s.writer.Write(p)
}

// NewStdoutHook creates a new stdout hook.
func NewStdoutHook(name string, processor *Processor, config []byte,
) (logrus.Hook, error) {
	c := &StdoutConfig{}
	if err := json.Unmarshal(config, c); err != nil {
		return nil, err
	}

	var logLevels []logrus.Level
	switch {
	case c.Level != "":
		logLevels = aboveLevel(c.Level)
	case len(c.Levels) > 0:
		logLevels = parseLevels(c.Levels)
	default:
		logLevels = logrus.AllLevels
	}

	w := writer.Hook{
		Writer: &stdoutWriter{
			processor: processor,
			writer:    stdout,
		},
		LogLevels: logLevels,
	}
	return &StdoutHook{w, processor}, nil
}
