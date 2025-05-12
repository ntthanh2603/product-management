package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}
func NewLogger() {
	encoder := getEncoderLog()
	core := zapcore.NewCore(
		encoder, 
		hook := lumber
		, zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
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