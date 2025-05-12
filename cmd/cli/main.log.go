package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Chạy lệnh sau để cài zap:
// go get -u go.uber.org/zap
func main() {

	// 1.
	// suger := zap.NewExample().Sugar()
	// suger.Infof("Hello, World") // {"level":"info","msg":"Hello, World"}

	// logger := zap.NewExample()
	// logger.Info("Hello, World", zap.String("name", "TuanThanh"), zap.Int("age", 25)) // {"level":"info","msg":"Hello, World","name":"TuanThanh","age":25}

	// 2.
	// Development
	// logger, _ := zap.NewDevelopment()
	// logger.Info("Development Logs") // 2025-03-29T16:12:26.629+0700    INFO    cli/main.log.go:16      Development Logs

	// // Production
	// logger, _ = zap.NewProduction()
	// logger.Info("Production Logs") // {"level":"info","ts":1743239546.6298692,"caller":"cli/main.log.go:20","msg":"Production Logs"}

	// 3.
	encoder := getEncoderLog()
	sync := getWriterSync()
	core := zapcore.NewCore(encoder, sync, zap.InfoLevel)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // Đảm bảo log được ghi trước khi chương trình kết thúc

	logger.Info("Info Logs", zap.Int("Line", 1))
	logger.Error("Error Logs", zap.Int("Line", 2))
}

// Format log
func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	// 1743239546.6298692 -> 2025-03-29T16:12:26.629+0700
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// ts -> Time
	encodeConfig.TimeKey = "time"

	// From info INFO
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// "caller":"cli/main.log.go:20"
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewConsoleEncoder(encodeConfig)
}

func getWriterSync() zapcore.WriteSyncer {
	// Đảm bảo thư mục .log tồn tại
	os.MkdirAll(".log", os.ModePerm)

	file, err := os.OpenFile(".log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(file), zapcore.AddSync(os.Stderr))
}
