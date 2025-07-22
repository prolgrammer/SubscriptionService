package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"subscription_service/config/pg"
)

type (
	Config struct {
		App  `mapstructure:"app"`
		HTTP `mapstructure:"http"`
		PG   pg.Config `mapstructure:"postgres"`
	}

	App struct {
		Name     string `mapstructure:"name"`
		Version  string `mapstructure:"version"`
		LogLevel string `mapstructure:"log_level"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

func New() (*Config, error) {
	cfg := Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	for _, key := range v.AllKeys() {
		anyVal := v.Get(key)
		str, ok := anyVal.(string)
		if !ok {
			continue
		}

		replaced := os.ExpandEnv(str)
		v.Set(key, replaced)
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	return &cfg, nil
}
