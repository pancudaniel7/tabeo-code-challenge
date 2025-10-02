package usecase

import (
	"context"
	"tabeo.org/challenge/internal/core/entity"
	"time"
)

type AppointmentUseCase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, a *entity.Appointment) error
	FindByVistDate(ctx context.Context, appointmentDate time.Time) (*entity.Appointment, error)
}

type HolidayHttpClient interface {
	RetrievePublicHolidays(year int, country string) ([]entity.PublicHolidays, error)
}

type HolidayCacheClient interface {
	GetPublicHolidays(ctx context.Context, year int, country string) ([]entity.PublicHolidays, error)
	SetPublicHolidays(ctx context.Context, year int, country string, holidays []entity.PublicHolidays) error
}
