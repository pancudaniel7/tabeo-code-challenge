package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"tabeo.org/challenge/internal/adapter/cache"
	"tabeo.org/challenge/internal/adapter/http"
	"tabeo.org/challenge/internal/adapter/repo"
	"tabeo.org/challenge/internal/core/usecase"
	"tabeo.org/challenge/internal/pkg/logger"

	"tabeo.org/challenge/internal/infra"
)

var (
	log logger.AppLogger
	db  *gorm.DB

	appointmentRepository usecase.AppointmentRepository
	holidayHttpClient     usecase.HolidayHttpClient
	holidaysCacheClient   usecase.HolidayCacheClient

	appointmentUsecase     usecase.AppointmentUseCase
	appointmentHttpHandler http.AppointmentHandler
)

func initComponents() {
	appointmentRepository = repo.NewAppointmentDefaultRepository(db)
	holidayHttpClient = http.NewHolidayClient()
	holidaysCacheClient = cache.NewHolidayCacheClient()
	appointmentUsecase = usecase.NewAppointmentDefaultUseCase(appointmentRepository, holidaysCacheClient, holidayHttpClient, log)
	appointmentHttpHandler = http.NewAppointmentDefaultHandler(log, appointmentUsecase)
}

func initRoutes(app *fiber.App) {
	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Up!")
	})
	app.Post("/appointments", appointmentHttpHandler.CreateAppointment)
}

func main() {
	log = infra.InitDefaultLogger()
	infra.InitDefaultConfig()
	db = infra.InitDB()

	initComponents()

	port := viper.GetInt("server.port")
	app := fiber.New()
	initRoutes(app)

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(context.Background(), err, fmt.Sprintf("Fatal error starting server: %s", err))
	}
}
