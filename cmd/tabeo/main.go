package main

import (
	"context"

	"tabeo.org/challenge/internal/infra"
	"tabeo.org/challenge/pkg/logger"
)

var (
	log logger.AppLogger
)

func main() {
	log = infra.InitDefaultLogger()
	infra.InitDefaultConfig()
	srv := infra.InitDefaultServer()

	log.Info(nil, "Starting server...", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(context.Background(), err, "Failed to start server")
	}
}
