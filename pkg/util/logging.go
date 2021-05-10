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

func MakeLogger(level string, formatter string) (logrus.FieldLogger, error) {
	log := logrus.New()
	var err error
	if log.Level, err = getLogrusLevel(level); err != nil {
		return nil, fmt.Errorf("setting log level failed: %w", err)
	}
	if log.Formatter, err = getLogrusFormatter(formatter); err != nil {
		return nil, fmt.Errorf("setting log formatter failed: %w", err)
	}

	return log, nil
}

func getLogrusLevel(level string) (logrus.Level, error) {
	res, ok := logrusLevels[level]
	if !ok {
		return 0, fmt.Errorf("%q is not a valid log level", level)
	}
	return res, nil
}

func getLogrusFormatter(typ string) (logrus.Formatter, error) {
	switch typ {
	case "text":
		return &logrus.TextFormatter{}, nil
	case "json":
		return &logrus.JSONFormatter{}, nil
	}
	return nil, fmt.Errorf("%q is not a valid log formatter", typ)
}
