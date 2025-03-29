package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database []struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		DBName   string `mapstructure:"dbName"`
	} `mapstructure:"database"`
}

func main() {
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
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println("Server Port::", config.Server.Port)

	for _, db := range config.Database {
		fmt.Printf("Database User: %s Password: %s, Host: %s, DBName: %s\n", db.User, db.Password, db.Host, db.DBName)
	}
}
