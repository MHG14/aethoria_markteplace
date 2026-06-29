package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPCfg     `mapstructure:"http"`
	Database PostgresCfg `mapstructure:"database"`
}

type HTTPCfg struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type PostgresCfg struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	fmt.Println("Loaded config:", v.ConfigFileUsed())

	return &cfg, nil
}
