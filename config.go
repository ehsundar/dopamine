package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("dopamine")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	viper.SetDefault("dsn", "./dopamine.sqlite3")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warning("config file not found. creating default config file")
			if err = viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
