package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)


type Config struct {
	Environment           string        `mapstructure:"ENVIRONMENT"`
	ServerAddress         string        `mapstructure:"SERVER_ADDRESS"`
	DatabaseURL           string        `mapstructure:"DATABASE_URL"`
	RedisURL              string        `mapstructure:"REDIS_URL"`
	TokenSymmetricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RateLimitDefault      int           `mapstructure:"RATE_LIMIT_DEFAULT"`
	RateLimitAuth         int           `mapstructure:"RATE_LIMIT_AUTH"`
	RateLimitEnabled      bool          `mapstructure:"RATE_LIMIT_ENABLED"`
	CORSAllowedOrigins    string        `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedMethods    string        `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowedHeaders    string        `mapstructure:"CORS_ALLOWED_HEADERS"`
	CORSAllowCredentials  bool          `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	LiveEnabled           bool          `mapstructure:"LIVE_ENABLED"`
	LiveUseMemoryBroker   bool          `mapstructure:"LIVE_USE_MEMORY_BROKER"`
}


func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")  
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "15m")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "24h")
	viper.SetDefault("RATE_LIMIT_DEFAULT", 100)
	viper.SetDefault("RATE_LIMIT_AUTH", 5)
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	
	
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	viper.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	viper.SetDefault("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Request-ID")
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)
	
	viper.SetDefault("LIVE_ENABLED", true)
	viper.SetDefault("LIVE_USE_MEMORY_BROKER", false) 

	err = viper.ReadInConfig()
	if err != nil {
		
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	
	if config.DatabaseURL == "" {
		err = fmt.Errorf("DATABASE_URL is required")
		return
	}

	if config.TokenSymmetricKey == "" {
		err = fmt.Errorf("TOKEN_SYMMETRIC_KEY is required (must be at least 32 characters)")
		return
	}

	if len(config.TokenSymmetricKey) < 32 {
		err = fmt.Errorf("TOKEN_SYMMETRIC_KEY must be at least 32 characters long")
		return
	}

	return
}
