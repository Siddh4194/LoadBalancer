package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {

	file, _ := os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	writeSyncer := zapcore.AddSync(file)

	encoder := zapcore.NewJSONEncoder(
		zap.NewProductionEncoderConfig(),
	)

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zapcore.InfoLevel,
	)

	Log = zap.New(core)
}
