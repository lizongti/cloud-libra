package log

import (
	"io"
	"sync"

	"github.com/mattn/go-colorable"
)

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
