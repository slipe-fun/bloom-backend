package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Driver   string
		Host     string
		Port     int
		User     string
		Password string
		DBName   string `mapstructure:"dbname"`
	}
	JWT struct {
		Secret      string
		ExpireHours int `mapstructure:"expire_hours"`
	}
	Argon2 struct {
		Time    uint32
		Memory  uint32
		Threads uint8
		KeyLen  uint32 `mapstructure:"key_len"`
	}
}

func LoadConfig(path string) *Config {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return &cfg
}

func (c *Config) JWTExpireDuration() time.Duration {
	return time.Duration(c.JWT.ExpireHours) * time.Hour
}
