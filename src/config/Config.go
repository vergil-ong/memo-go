package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitViperConfig() {
	wd, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wd + "/config")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fail to read file %s \n", err))
	}
}
