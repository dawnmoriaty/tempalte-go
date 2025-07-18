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

type Config struct {
	// tránh việc viper không nhận diện được cấu trúc lồng nhau
	// sử dụng `squash` để gom các trường con vào cấu trúc cha
	Database DatabaseConfig `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
}

var (
	cfg *Config
)

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Error unmarshalling config: " + err.Error())
	}

	if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.DBName == "" || cfg.Database.HTTPPort == "" || cfg.Database.Port == "" {
		panic("Database configuration is incomplete. Please check environment variables or .env file.")
	}

	return cfg
}

func GetConfig() *Config {
	if cfg == nil {
        LoadConfig() 
    }
	return cfg
}
