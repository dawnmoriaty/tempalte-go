package configs

type RedisConfig struct {
    Addr     string `mapstructure:"REDIS_ADDR"`
    Password string `mapstructure:"REDIS_PASSWORD"`
    DB       int    `mapstructure:"REDIS_DB"`
}