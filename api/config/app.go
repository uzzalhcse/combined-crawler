package config

import "github.com/spf13/viper"

type AppConfig struct {
	Name      string
	Env       string
	Debug     bool
	Port      string
	JwtSecret string
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Name:      viper.GetString("APP_NAME"),
		Port:      viper.GetString("APP_PORT"),
		JwtSecret: viper.GetString("JWT_SECRET"),
	}
}
