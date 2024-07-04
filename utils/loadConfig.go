package utils

import (
	"github.com/spf13/viper"
)


type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// incase i have a env file 

	viper.AutomaticEnv() 

	err  = viper.ReadInConfig();

	if err != nil {
		return
	}

	// go ahead and unmarshall the config

	err = viper.Unmarshal(&config)
	return

}
