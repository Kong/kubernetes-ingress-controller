package util

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
	"panic": zap.ErrorLevel,
	"fatal": zap.ErrorLevel,
	"error": zap.ErrorLevel,
	"warn":  zap.InfoLevel,
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

	encoder, err := getZapEncoding(formatter)
	if err != nil {
		return nil, fmt.Errorf("setting log formatter failed: %w", err)
	}
	// TODO 1893 we can maybe avoid building from a core by using
	// https://pkg.go.dev/go.uber.org/zap#RegisterSink to make a URL
	// that points to an io.Writer somehow. Not clear it'd be worth it.
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

func getZapEncoding(typ string) (zapcore.Encoder, error) {
	switch typ {
	case "text", "console":
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       zapcore.OmitKey,
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}), nil
	case "json":
		return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       zapcore.OmitKey,
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}), nil
	}
	return nil, fmt.Errorf("%q is not a valid log formatter", typ)
}

// for reference, these are the standard zap production and development configurations

var (
	prod = zap.Config{ //nolint:unused
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	devel = zap.Config{ //nolint:unused
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
)
