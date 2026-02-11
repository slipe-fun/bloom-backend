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
	Email struct {
		Email        string
		SmtpHost     string `mapstructure:"smtp_host"`
		SmtpPort     string `mapstructure:"smtp_port"`
		SmtpLogin    string `mapstructure:"smtp_login"`
		SmtpPassword string `mapstructure:"smtp_password"`
	}
	GoogleAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
		BundleID     string `mapstructure:"bundle_id"`
	} `mapstructure:"google_auth"`
	RateLimit struct {
		Enabled                  bool `mapstructure:"enabled"`
		AuthRequestsPerMinute    int  `mapstructure:"auth_requests_per_minute"`
		GeneralRequestsPerMinute int  `mapstructure:"general_requests_per_minute"`
		WindowMinutes            int  `mapstructure:"window_minutes"`
	} `mapstructure:"rate_limit"`
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
