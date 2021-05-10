package util

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	logrusLevels = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}
	logrusFormats = map[string]logrus.Formatter{
		"text": &logrus.TextFormatter{},
		"json": &logrus.JSONFormatter{},
	}
)

func GetLogrusLevel(level string) (logrus.Level, error) {
	res, ok := logrusLevels[level]
	if !ok {
		return 0, fmt.Errorf("%q is not a valid log level", level)
	}
	return res, nil
}

func GetLogrusFormatter(typ string) (logrus.Formatter, error) {
	switch typ {
	case "text":
		return &logrus.TextFormatter{}, nil
	case "json":
		return &logrus.JSONFormatter{}, nil
	}
	return nil, fmt.Errorf("%q is not a valid log formatter", typ)
}
