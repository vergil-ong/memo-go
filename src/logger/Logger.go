package logger

import (
	"encoding/json"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var Logger *zap.Logger
var once sync.Once

func InitLogger() *zap.Logger {
	once.Do(func() {
		writer := getLogWriter()
		encoder := getEncoder()
		core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

		//AddCaller 显示文件和行号
		Logger = zap.New(core, zap.AddCaller())
	})

	return Logger
}

func getLogWriter() zapcore.WriteSyncer {
	logFilePath := viper.GetString("log.file-path")
	logger := lumberjack.Logger{
		Filename:   logFilePath,
		MaxAge:     10,
		MaxSize:    5,
		MaxBackups: 3,
		Compress:   false,
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(&logger), zapcore.AddSync(os.Stdout))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func GetJson(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
