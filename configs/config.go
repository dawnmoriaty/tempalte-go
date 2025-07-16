package configs

import (
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Port     string `mapstructure:"DB_PORT"`
	DBName   string `mapstructure:"DB_NAME"`
	HTTPPort string `mapstructure:"HTTP_PORT"`
}

var (
	cfg *DatabaseConfig
)

func LoadConfig() *DatabaseConfig {
	viper.SetConfigFile(".env")
	viper.ReadInConfig() 

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Error unmarshalling config: " + err.Error())
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" || cfg.HTTPPort == "" || cfg.Port == "" {
		panic("Database configuration is incomplete. Please check environment variables or .env file.")
	}

	return cfg
}

func GetConfig() *DatabaseConfig {
	return cfg
}
