// +build windows nacl plan9

package hook

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type hookCreaterFunc func(string, *Processor, []byte) (logrus.Hook, error)

var hooks = map[string]hookCreaterFunc{
	"lumberjack":     NewLumberjackHook,
	"lumberjack_err": NewLumberjackErrHook,
	"stdout":         NewStdoutHook,
	"stderr":         NewStderrHook,
}

// New return a hook init by yaml config
func New(name string, typ string, processor *Processor,
	config []byte) (logrus.Hook, error) {
	creature, ok := hooks[typ]
	if !ok {
		return nil, fmt.Errorf("no hook %s is found", typ)
	}
	return creature(name, processor, config)
}
