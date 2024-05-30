package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the entity responsible for logging.
type Logger struct {
	*zap.Logger
}

// New returns a new instance of a logger.
func New(logLevel string) *Logger {
	cfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		TimeKey:        "time",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), os.Stdout, level)

	return &Logger{zap.New(core)}
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}
