package formatter

import (
	"errors"
	"fmt"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var ErrUnexpectedFieldKey = errors.New("unexpected field key")

type JSONFormatter struct {
	logrus.JSONFormatter
}

func (*Generator) JSON(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
	fieldMap := make(logrus.FieldMap)

	for k, v := range h.GetStringMapString("fieldMap") {
		switch k {
		case logrus.FieldKeyMsg:
			fieldMap[logrus.FieldKeyMsg] = v
		case logrus.FieldKeyLevel:
			fieldMap[logrus.FieldKeyLevel] = v
		case logrus.FieldKeyTime:
			fieldMap[logrus.FieldKeyTime] = v
		case logrus.FieldKeyLogrusError:
			fieldMap[logrus.FieldKeyLogrusError] = v
		case logrus.FieldKeyFunc:
			fieldMap[logrus.FieldKeyFunc] = v
		case logrus.FieldKeyFile:
			fieldMap[logrus.FieldKeyFile] = v
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnexpectedFieldKey, k)
		}
	}

	formatter := &JSONFormatter{
		JSONFormatter: logrus.JSONFormatter{
			DisableTimestamp:  true,
			DisableHTMLEscape: true,
			FieldMap:          fieldMap,
			PrettyPrint:       false,
		},
	}

	return formatter, nil
}
