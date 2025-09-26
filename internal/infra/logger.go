package infra

import (
	"github.com/spf13/viper"
	"tabeo.org/challenge/pkg/logger"
)

// InitDefaultLogger initializes and returns a AppLogger using the provided Config.
func InitDefaultLogger() logger.AppLogger {
	cfg := logger.Config{
		Level: viper.GetString("log.level"),
		JSON:  viper.GetBool("log.json"),
	}
	return logger.NewDefaultLogger(cfg)
}
