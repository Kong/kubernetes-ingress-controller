package util

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// we currently implement two different loggers and use a middleware called
// logrusr to translate logrus loggers into go-logrs (used by controller-runtime).
// The middleware currently squashes loglevels 0-4 together and effectively starts
// the "info" level logging at 0 (whereas logrus starts that at 4 (range from 0)).
// Since the middleware at the time of writing made all of this part of the private
// implementation these options are used for convenience until a time when we unify
// our logging implementations into one or the other.
//
// See: https://github.com/Kong/kubernetes-ingress-controller/issues/1893
const (
	logrusrDiff = 4

	// InfoLevel is the converted logging level from logrus to go-logr for
	// information level logging. Note that the logrusr middleware technically
	// flattens all levels prior to this level into this level as well.
	InfoLevel = int(logrus.InfoLevel) - logrusrDiff

	// DebugLevel is the converted logging level from logrus to go-logr for
	// debug level logging.
	DebugLevel = int(logrus.DebugLevel) - logrusrDiff

	// WarnLevel is the converted logrus level to go-logr for warnings
	WarnLevel = int(logrus.WarnLevel) - logrusrDiff
)

var logrusLevels = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"trace": logrus.TraceLevel,
}

func MakeLogger(level string, formatter string) (logrus.FieldLogger, error) {
	log := logrus.New()
	var err error

	logLevel, err := getLogrusLevel(level)
	if err != nil {
		return nil, fmt.Errorf("setting log level failed: %w", err)
	}
	if log.Formatter, err = getLogrusFormatter(formatter); err != nil {
		return nil, fmt.Errorf("setting log formatter failed: %w", err)
	}

	log.SetLevel(logLevel)
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
