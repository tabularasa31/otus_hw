package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct { // TODO
}

func New(level string) (*zap.SugaredLogger, error) {
	loglevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return &zap.SugaredLogger{}, fmt.Errorf("failed to parse log level: %v", err)
	}

	writerSyncer := zapcore.AddSync(os.Stdout)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(encoder, writerSyncer, loglevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	return sugarLogger, nil
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
