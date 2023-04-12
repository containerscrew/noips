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
	LogLevel string // by the moment only supported info or debu
}

func NewLogger(options *LoggerOptions) *CustomLogger {
	return &CustomLogger{Log: logger(options)}
}

func logger(options *LoggerOptions) *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = append(config.OutputPaths, "logs.log")
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	defaultLogLevel, _ := zapcore.ParseLevel(options.LogLevel)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	customLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return customLogger.Sugar()
}
