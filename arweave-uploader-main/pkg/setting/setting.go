package setting

import (
	"ccian.cc/really/arweave-api/pkg/util"
)

var gConfig Config

func Setup(file string) error {
	return util.ReadConfig(file, &gConfig, nil)
}

func AllConfig() *Config {
	return &gConfig
}

func App() *AppConfig {
	return &gConfig.App
}

func Logger() *LoggerConfig {
	return &gConfig.Logger
}

func Server() *ServerConfig {
	return &gConfig.Server
}
