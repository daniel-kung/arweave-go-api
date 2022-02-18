package util

import (
	"strings"

	"github.com/spf13/viper"
)

func ReadConfig(filePath string, out interface{}, defaults map[string]interface{}) error {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetConfigFile(filePath)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // for nested structure
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&out); err != nil {
		return err
	}

	return nil
}
