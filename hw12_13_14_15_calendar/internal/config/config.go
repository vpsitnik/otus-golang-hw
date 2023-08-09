package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger *LoggerConf
	Db     *DatabaseConf
}

type LoggerConf struct {
	Level string
}

type DatabaseConf struct {
	Type     string
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func NewConfig(configFile string) *Config {
	viper.SetDefault("logger.level", "INFO")
	viper.SetDefault("database.type", "in-memory")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &Config{
		Logger: &LoggerConf{
			Level: viper.GetString("logger.level"),
		},
		Db: &DatabaseConf{
			Type:     viper.GetString("database.type"),
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			Database: viper.GetString("database.database"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
		},
	}
}

// TODO
