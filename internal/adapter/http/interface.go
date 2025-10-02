package http

import (
	"github.com/gofiber/fiber/v3"
	"tabeo.org/challenge/internal/core/entity"
)

type AppointmentHandler interface {
	CreateAppointment(ctx fiber.Ctx) (*AppointmentResponse, error)
}

type HolidayClient interface {
	RetrievePublicHolidays(year int, country string) ([]entity.PublicHolidays, error)
}
