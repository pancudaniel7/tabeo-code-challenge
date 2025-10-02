package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"tabeo.org/challenge/internal/core/usecase"
	"tabeo.org/challenge/internal/pkg/apperr"
	"tabeo.org/challenge/internal/pkg/logger"
)

var validate = validator.New()

type AppointmentDefaultHandler struct {
	log                logger.AppLogger
	appointmentUseCase usecase.AppointmentUseCase
}

func NewAppointmentDefaultHandler(log logger.AppLogger, appointmentUseCase usecase.AppointmentUseCase) AppointmentHandler {
	return &AppointmentDefaultHandler{log: log, appointmentUseCase: appointmentUseCase}
}

func (b *AppointmentDefaultHandler) CreateAppointment(ctx fiber.Ctx) error {
	var req AppointmentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	appointmentDTO, err := req.ToEntity()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	appointmentEntity, err := b.appointmentUseCase.CreateAppointment(ctx, appointmentDTO)
	if err != nil {
		return apperr.HttpHandleError(ctx, err)
	}
	
	res := req.ToDTO(appointmentEntity)
	return ctx.Status(fiber.StatusCreated).JSON(res)
}
