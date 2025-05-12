package initialize

import (
	"fmt"
	"myproject/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./config") // path config
	viper.SetConfigName("local")    // tên file
	viper.SetConfigType("yaml")     // điều hiện file

	// read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// read server config not use config struct
	// fmt.Println("Server Port::", viper.GetInt("server.port"))

	// config struct
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}