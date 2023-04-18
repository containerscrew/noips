package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomLogger struct {
	Log *zap.SugaredLogger
}

type LoggerOptions struct {
	LogLevel string // by the moment only supported info or debug
}

func NewLogger(options *LoggerOptions) *CustomLogger {
	return &CustomLogger{Log: logger(options)}
}

func logger(options *LoggerOptions) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "timestamp"
	//config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewJSONEncoder(config)

	defaultLogLevel, _ := zapcore.ParseLevel(options.LogLevel)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	customLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return customLogger.Sugar()
}
