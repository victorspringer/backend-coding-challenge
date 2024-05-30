package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// TODO: add comments

type Config struct {
	UserService struct {
		Environment string `mapstructure:"environment"`
		LogLevel    string `mapstructure:"log_level"`
		Server      struct {
			Port         string        `mapstructure:"port"`
			ReadTimeout  time.Duration `mapstructure:"read_timeout"`
			WriteTimeout time.Duration `mapstructure:"write_timeout"`
			IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		} `mapstructure:"server"`
	} `mapstructure:"user_service"`
}

func New(env string) (*Config, error) {
	cfg := &Config{}
	viper.AutomaticEnv()
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	viper.SetEnvPrefix("User")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if env == "" {
		env = "development"
	}

	viper.SetConfigName(env)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
