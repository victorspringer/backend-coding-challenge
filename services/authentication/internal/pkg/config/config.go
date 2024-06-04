package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config contains the application configuration parameters.
type Config struct {
	AuthenticationService struct {
		Environment string `mapstructure:"environment"`
		LogLevel    string `mapstructure:"log_level"`
		PrivateKey  string `mapstructure:"private_key"`
		PublicKey   string `mapstructure:"public_key"`
		Claims      struct {
			Issuer                      string        `mapstructure:"issuer"`
			AccessTokenExpiration       time.Duration `mapstructure:"access_token_expiration"`
			AnonymousExpiration         time.Duration `mapstructure:"anonymous_expiration"`
			ShortRefreshTokenExpiration time.Duration `mapstructure:"short_refresh_token_expiration"`
			LongRefreshTokenExpiration  time.Duration `mapstructure:"long_refresh_token_expiration"`
		} `mapstructure:"claims"`
		Server struct {
			Port         string        `mapstructure:"port"`
			ReadTimeout  time.Duration `mapstructure:"read_timeout"`
			WriteTimeout time.Duration `mapstructure:"write_timeout"`
			IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		} `mapstructure:"server"`
	} `mapstructure:"authentication_service"`
	UserService struct {
		Endpoint string        `mapstructure:"endpoint"`
		Timeout  time.Duration `mapstructure:"timeout"`
	} `mapstructure:"user_service"`
	Redis struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"redis"`
}

// New returns a new instance of Config.
func New(env string) (*Config, error) {
	cfg := &Config{}
	viper.AutomaticEnv()
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	viper.SetEnvPrefix("AuthenticationService")
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
