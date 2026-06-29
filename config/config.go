package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
	DBName     string `mapstructure:"DBNAME"`
	SERVERPORT string `mapstructure:"SERVERPORT"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func Load() (Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
