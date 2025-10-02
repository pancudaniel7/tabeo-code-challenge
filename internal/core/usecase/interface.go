package usecase

import (
	"context"
	"tabeo.org/challenge/internal/core/entity"
)

type AppointmentUseCase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (string, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, a *entity.Appointment) error
}
