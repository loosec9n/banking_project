package utils

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct will hold all configuration od application
// the values are readby viper from config files or env variables
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("viper: error reading env files")
		return
	}

	err = viper.Unmarshal(&config)
	return
}
