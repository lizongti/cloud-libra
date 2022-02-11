package log

import (
	"io/ioutil"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
	"github.com/sirupsen/logrus"
)

func NewLogger(config *tree.Tree) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetReportCaller(true)

	formatter, err := NewFormatterCreater(cast.ToString(config.Get(magic.UnixChain("formatter"))))
	if err != nil {
		return nil, err
	}
	logger.SetFormatter(formatter)

	logger.SetOutput(ioutil.Discard)

	logLevel, err := NewLogLevel(cast.ToString(config.Get(magic.UnixChain("level"))))
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel.Level())

	for name, c := range cast.ToStringMap(config.Get(magic.UnixChain("hooks"))) {
		hook, err := NewHook(name, tree.NewMapTree().SetData(cast.ToStringMap(c)))
		if err != nil {
			return nil, err
		}
		logger.AddHook(hook)
	}

	return logger, nil
}
