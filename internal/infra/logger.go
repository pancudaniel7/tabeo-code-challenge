package infra

import (
	"github.com/spf13/viper"
	logger2 "tabeo.org/challenge/internal/pkg/logger"
)

// InitDefaultLogger initializes and returns a AppLogger using the provided Config.
func InitDefaultLogger() logger2.AppLogger {
	cfg := logger2.Config{
		Level: viper.GetString("log.level"),
		JSON:  viper.GetBool("log.json"),
	}
	return logger2.NewDefaultLogger(cfg)
}
