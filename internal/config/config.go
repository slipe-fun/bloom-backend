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
	Redis struct {
		Adress   string
		Password string
	}
	JWT struct {
		Secret      string
		ExpireHours int `mapstructure:"expire_hours"`
	}
	RateLimit struct {
		Enabled                  bool `mapstructure:"enabled"`
		AuthRequestsPerMinute    int  `mapstructure:"auth_requests_per_minute"`
		GeneralRequestsPerMinute int  `mapstructure:"general_requests_per_minute"`
		WindowMinutes            int  `mapstructure:"window_minutes"`
	} `mapstructure:"rate_limit"`
	WebAuthn struct {
		RPID          string   `mapstructure:"rpid"`
		RPDisplayName string   `mapstructure:"rp_display_name"`
		RPOrigins     []string `mapstructure:"rp_origins"`
	} `mapstructure:"webauthn"`
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

func (c *Config) RateLimitWindow() time.Duration {
	return time.Duration(c.RateLimit.WindowMinutes) * time.Minute
}
