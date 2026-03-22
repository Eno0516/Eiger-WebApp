package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env         string `mapstructure:"ENV"`
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
