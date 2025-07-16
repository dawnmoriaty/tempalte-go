package logger

import (
	"os"
	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

func InitLogger() {
	writer := getLogWriter()
	consoleWriter := zapcore.Lock(os.Stdout)
	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, zapcore.InfoLevel),        // Log ra file từ level Info trở lên
		zapcore.NewCore(encoder, consoleWriter, zapcore.DebugLevel), // Log ra console từ level Debug trở lên
	)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	zap.ReplaceGlobals(logger)
}

func GetLogger() *zap.SugaredLogger {
	return sugarLogger
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(config)
}

func getLogWriter() zapcore.WriteSyncer {		
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   false,
	}
	return zapcore.AddSync(lumberjackLogger)
}
