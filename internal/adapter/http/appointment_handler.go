package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"tabeo.org/challenge/internal/pkg/logger"
)

var validate = validator.New()

type AppointmentDefaultHandler struct {
	log logger.AppLogger
}

func NewAppointmentDefaultHandler(log logger.AppLogger) AppointmentHandler {
	return &AppointmentDefaultHandler{log: log}
}

func (b *AppointmentDefaultHandler) CreateAppointment(c fiber.Ctx) (*AppointmentResponse, error) {
	var req AppointmentRequest
	if err := c.Bind().Body(&req); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(&req); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	book, err := req.ToEntity()
	if err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO: Call the use case to create the appointment
	print(book)

	return nil, c.SendStatus(fiber.StatusNotImplemented)
}
