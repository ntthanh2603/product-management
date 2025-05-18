package initialize

import (
	"myproject/global"
	"myproject/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}