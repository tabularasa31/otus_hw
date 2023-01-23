package logger

import (
	"os"

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
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(encoder, writerSyncer, loglevel)

	sugar := zap.New(core).Sugar()

	return sugar
}
