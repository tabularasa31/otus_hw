package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	loglevel, err := zap.ParseAtomicLevel(level)

	if err != nil {
		return &Logger{
			logger: &zap.SugaredLogger{},
		}
	}

	writerSyncer := zapcore.AddSync(os.Stdout)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(encoder, writerSyncer, loglevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	return &Logger{
		logger: sugarLogger,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.Debug(message, args)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.Info(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.Warn(message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.Error(message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.Fatal(message, args...)

	os.Exit(1)
}
