package setting

type Config struct{
	Mysql MySQLSetting `mapstructure:"mysql"`
}

type MySQLSetting struct{
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	User string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`
	MaxIdleConns int `mapstructure:"maxIdleConns"`
	MaxOpenConns int `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int `mapstructure:"connMaxLifetime"`
}