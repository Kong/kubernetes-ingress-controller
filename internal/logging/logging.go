package logging

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logr does not use semantic log types. Logs are either informational or error messages with arbitrary verbosity,
// though errors print regardless of the verbosity (though ensuring this is left up to the implementation--logr
// interfaces only set the level; they do not filter lines). Only levels above the base info level have meaning.
// The error and warn levels here are effectively just placeholders.

const (
	// ErrorLevel is the logr verbosity level for errors.
	ErrorLevel = 0

	// WarnLevel is the logr verbosity level for warnings.
	WarnLevel = 0

	// InfoLevel is the logr verbosity level for info logs.
	InfoLevel = 0

	// DebugLevel is the logr verbosity level for debug logs.
	DebugLevel = 1

	// TraceLevel is the logr verbosity level for trace logs.
	TraceLevel = 2
)

// similar to the above, these are squashed for levels info and higher. assumptions by other libraries don't easily
// allow differentiation between info and warn.

var zapLevels = map[string]zapcore.Level{
	"error": zap.ErrorLevel,
	"info":  zap.InfoLevel,
	"debug": zap.DebugLevel,
	// zap has no stock trace level, but does accept it for filtering if you define your own
	"trace": zapcore.Level(-2),
}

func MakeLogger(level string, formatter string, output io.Writer) (*zap.Logger, error) {
	logLevel, err := getZapLevel(level)
	if err != nil {
		return nil, fmt.Errorf("setting log level failed: %w", err)
	}
	levelFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// note that zapr flips the sign of Info V-levels, so V(2) results in lvl=-2 here
		return lvl >= logLevel
	})

	encoder, err := GetZapEncoding(formatter)
	if err != nil {
		return nil, fmt.Errorf("setting log formatter failed: %w", err)
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(output), levelFunc)

	return zap.New(core), nil
}

func getZapLevel(level string) (zapcore.Level, error) {
	res, ok := zapLevels[level]
	if !ok {
		return 0, fmt.Errorf("%q is not a valid log level", level)
	}
	return res, nil
}

func GetZapEncoding(typ string) (zapcore.Encoder, error) {
	switch typ {
	case "text", "console":
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}), nil
	case "json":
		return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}), nil
	}
	return nil, fmt.Errorf("%q is not a valid log formatter", typ)
}
