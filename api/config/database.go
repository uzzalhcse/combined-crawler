package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func loadDBConfig() DatabaseConfig {
	return DatabaseConfig{
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Name:     viper.GetString("DB_DATABASE"),
	}
}
