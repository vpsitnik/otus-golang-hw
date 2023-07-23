package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger *LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(configFile string) *Config {
	viper.SetDefault("logger.level", "ERROR")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	log.Printf("Log level: %v\n", viper.GetString("logger.level"))
	return &Config{
		Logger: &LoggerConf{
			Level: viper.GetString("logger.level"),
		},
	}
}

// TODO
