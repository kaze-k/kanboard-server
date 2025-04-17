package global

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"server/internal/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getLogger() *zap.SugaredLogger {
	level := zapcore.DebugLevel
	if constant.EnvConfig.Mode == "debug" {
		level = zapcore.InfoLevel
	}

	consoleEncoder := getEncoder(true)
	fileEncoder := getEncoder(false)

	console := zapcore.Lock(zapcore.AddSync(os.Stdout))
	writer := zapcore.Lock(zapcore.AddSync(GetWriter()))

	consoleCore := zapcore.NewCore(consoleEncoder, console, level)
	fileCore := zapcore.NewCore(fileEncoder, writer, level)
	core := zapcore.NewTee(consoleCore, fileCore)

	return zap.New(core).Sugar()
}

func getEncoder(color bool) zapcore.Encoder {
	encodeLevel := zapcore.CapitalLevelEncoder
	if color {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Local().Format(time.DateTime))
		},
		EncodeLevel:  encodeLevel,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	if color {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func GetWriter() *lumberjack.Logger {
	separator := string(filepath.Separator)
	logFilePath := fmt.Sprintf(
		"%s%s%s%s%s.log",
		constant.LogConfig.Path,
		separator,
		constant.ServerConfig.Name,
		separator,
		time.Now().Format(time.DateOnly),
	)

	writer := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    constant.LogConfig.MaxSize,
		MaxBackups: constant.LogConfig.MaxBackups,
		MaxAge:     constant.LogConfig.MaxAge,
		Compress:   constant.LogConfig.Compress,
	}

	return writer
}
