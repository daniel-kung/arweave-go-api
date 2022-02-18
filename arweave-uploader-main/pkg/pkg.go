package pkg

import (
	"ccian.cc/really/arweave-api/pkg/logging"
	"ccian.cc/really/arweave-api/pkg/setting"
)

func Setup(configFile string) error {
	if err := setting.Setup(configFile); err != nil {
		return err
	}
	logging.Setup(setting.Logger())
	return nil
}
