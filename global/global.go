package global

import (
	"myproject/pkg/logger"
	"myproject/pkg/setting"

	"gorm.io/gorm"
)

var (
	Config setting.Config
	Logger *logger.LoggerZap
	Mdb *gorm.DB
)

