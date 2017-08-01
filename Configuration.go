// Config.go
package autositespeed

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(configPathList []string, configName string, conf interface{}) {

	for index := 0; index < len(configPathList); index++ {
		viper.AddConfigPath(configPathList[index])
	}

	viper.SetConfigName(configName)

	viper.ReadInConfig()

	viper.Unmarshal(conf)

	fmt.Println(conf)
}
