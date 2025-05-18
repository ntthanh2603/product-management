package initialize

import (
	"fmt"
	"myproject/global"

	"go.uber.org/zap"
)

func Run(){
	LoadConfig()
	m := global.Config.Mysql
	fmt.Println("Load config mysql:", m.User, m.Password)
	InitLogger()
	global.Logger.Info("Load config logger", zap.String("ok", "success"))
	InitRedis()
	InitPostgres()

	r := InitRouter()

	r.Run(":8080")
}