package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
	DBName     string `mapstructure:"DBNAME"`
	SERVERPORT string `mapstructure:"SERVERPORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

func Load() (Config, error) {
	var cfg Config

	// همیشه ENV فعال باشه
	viper.AutomaticEnv()

	// فقط اگر فایل واقعاً وجود داشت بخونش
	if _, err := os.Stat("./config/config.env"); err == nil {
		viper.SetConfigFile("./config/config.env")

		if err := viper.ReadInConfig(); err != nil {
			log.Println("config file exists but failed to read:", err)
		} else {
			log.Println("config file loaded")
		}
	} else {
		log.Println("no config file found → using ENV (Render mode)")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
