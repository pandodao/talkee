package config

import (
	"github.com/fox-one/pkg/config"
)

var cfg *Config

func Load(cfgFile string, cfg *Config) error {
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}
	return nil
}

func InitConfig() {
	cfg = new(Config)
	if err := Load("./config.yaml", cfg); err != nil {
		panic(err)
	}
}

func C() *Config {
	if cfg == nil {
		InitConfig()
	}
	return cfg
}
