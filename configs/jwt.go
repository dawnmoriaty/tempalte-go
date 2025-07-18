package configs

type JWTConfig struct {
	AccessTokenLife     string `mapstructure:"JWT_ACCESS_TOKEN_LIFE"`
	RefreshTokenLife    string `mapstructure:"JWT_REFRESH_TOKEN_LIFE"`
	AccessTokenSecret   string `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret  string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
}