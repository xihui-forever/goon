package config

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/spf13/viper"
	"os"
)

func Load() {
	//log.SetOutput(log.GetOutputWriterHourly(filepath.Join(utils.GetExecPath(), "goon"), 12))
	log.SetLevel(log.InfoLevel)

	viper.AddConfigPath(utils.GetExecPath())
	viper.AddConfigPath(utils.GetPwd())

	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// 配置一些默认值

	err := viper.ReadInConfig()
	if err != nil {
		switch e := err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Debug("not found conf file, use default")
		case *os.PathError:
			log.Debug("not find conf file in %s", e.Path)
		default:
			log.Fatalf("load config fail:%v", err)
			return
		}
	}

	//viper.WatchConfig()
}
