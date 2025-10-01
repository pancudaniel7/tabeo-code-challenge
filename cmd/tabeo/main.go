package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"tabeo.org/challenge/internal/pkg/logger"

	"tabeo.org/challenge/internal/infra"
)

var (
	log logger.AppLogger
)

func main() {
	log = infra.InitDefaultLogger()
	infra.InitDefaultConfig()

	port := viper.GetInt("server.port")
	app := fiber.New()
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(context.Background(), err, fmt.Sprintf("Fatal error starting server: %s", err))
	}
}
