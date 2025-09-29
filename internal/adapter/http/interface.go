package http

import "github.com/gofiber/fiber/v3"

type AppointmentHandler interface {
	CreateAppointment(ctx fiber.Ctx) (*AppointmentResponse, error)
}
