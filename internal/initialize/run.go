package initialize

import (
	"fmt"
	"myproject/global"
)

func Run(){
	LoadConfig()
	m := global.Config.Mysql
	fmt.Println("Load config mysql:", m.User, m.Password)
	InitLogger()
	InitRedis()
	InitPostgres()

	r := InitRouter()

	r.Run(":8080")
}