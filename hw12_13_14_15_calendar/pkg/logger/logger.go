package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.SugaredLogger {
	loglevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		atom := zap.NewAtomicLevel()
		atom.SetLevel(zap.ErrorLevel)
		loglevel = atom
	}

	writerSyncer := zapcore.AddSync(os.Stdout)

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Layout)
	encoder := zapcore.NewJSONEncoder(loggerConfig.EncoderConfig)

	core := zapcore.NewCore(encoder, writerSyncer, loglevel)

	sugar := zap.New(core).Sugar()

	return sugar
}
