package infra

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3/log"
	"github.com/spf13/viper"
)

// InitDefaultConfig initializes the viper configuration
func InitDefaultConfig() {
	viper.SetEnvPrefix("TABEO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.BindEnv("CONFIG_NAME"); err != nil {
		log.Fatal(fmt.Sprintf("Fatal error binding env variable: %s", err))
	}

	viper.AutomaticEnv()
	configName := viper.GetString("CONFIG_NAME")
	if configName == "" {
		configName = "local"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("Fatal error reading config file: %s", err))
	}
}
