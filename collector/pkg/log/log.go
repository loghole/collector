package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level   string
	Options []Option
}

type Option func(options []zap.Option) []zap.Option

func AddCaller() Option {
	return func(options []zap.Option) []zap.Option {
		return append(options, zap.AddCaller())
	}
}

func AddStacktrace(level string) Option {
	return func(options []zap.Option) []zap.Option {
		return append(options, zap.AddStacktrace(parseLevel(level)))
	}
}

func WithField(key string, value interface{}) Option {
	return func(options []zap.Option) []zap.Option {
		return append(options, zap.Fields(zap.Any(key, value)))
	}
}

func NewLogger(config *Config) (*zap.SugaredLogger, error) {
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder(), os.Stdout, parseLevel(config.Level)),
	}

	options := make([]zap.Option, 0, len(config.Options))

	for _, option := range config.Options {
		options = option(options)
	}

	return zap.New(zapcore.NewTee(cores...), options...).Sugar(), nil
}

func parseLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "err", "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func consoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}
