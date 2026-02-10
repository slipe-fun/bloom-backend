package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/google/uuid"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func Init(logsPath string) error {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logsPath,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), consoleWriter, zapcore.DebugLevel),
	)

	l := zap.New(core, zap.AddCaller())

	Log = l
	return nil
}

func LogError(msg string, service string) {
	requestID := uuid.New().String()

	Log.Error(
		msg,
		zap.String("service", service),
		zap.String("request_id", requestID),
	)
}
