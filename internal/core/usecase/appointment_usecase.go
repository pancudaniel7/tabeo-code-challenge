package usecase

import (
	"context"
	"github.com/spf13/viper"
	"tabeo.org/challenge/internal/adapter/cache"
	"tabeo.org/challenge/internal/adapter/http"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
	"tabeo.org/challenge/internal/pkg/logger"
	"time"
)

type AppointmentDefaultUseCase struct {
	appointmentRepo     AppointmentRepository
	holidaysCacheClient cache.HolidayCacheClient
	holidaysHttpClient  http.HolidayClient
	logger              logger.AppLogger
}

func NewAppointmentDefaultUseCase(appointmentRepo AppointmentRepository, cacheClient cache.HolidayCacheClient, holidayClient http.HolidayClient, logger logger.AppLogger) AppointmentUseCase {
	return &AppointmentDefaultUseCase{
		appointmentRepo:     appointmentRepo,
		holidaysCacheClient: cacheClient,
		holidaysHttpClient:  holidayClient,
		logger:              logger,
	}
}

// CreateAppointment creates a new appointment and checks for public holidays on the visit date.
func (a AppointmentDefaultUseCase) CreateAppointment(ctx context.Context, appointment *entity.Appointment) (string, error) {
	dismiss, err := a.dismissAppointment(ctx, appointment)
	if err != nil {
		return "", err
	}
	if dismiss {
		return "", apperr.InvalidArgument("appointment date falls on a public holiday", nil)
	}
	if err := a.appointmentRepo.Create(ctx, appointment); err != nil {
		return "", apperr.Internal("failed to create appointment", err)
	}

	return appointment.ID.String(), nil
}

// dismissAppointment dismiss appointment because the date is a public holiday
func (a AppointmentDefaultUseCase) dismissAppointment(ctx context.Context, appointment *entity.Appointment) (bool, error) {
	year := appointment.VisitDate.Year()
	country := viper.GetString("holidays.country")
	publicHolidays, err := a.retrieverPublicHolidays(ctx, year, country)
	if err != nil {
		a.logger.Error(ctx, err, "failed to retrieve public holidays during appointment creation")
	}

	for _, holiday := range publicHolidays {
		holidayDate, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			a.logger.Error(ctx, err, "failed to parse holiday date", "holiday", holiday.Name)
			continue
		}
		if holidayDate.Equal(appointment.VisitDate.Truncate(24 * time.Hour)) {
			a.logger.Warn(ctx, "appointment falls on a public holiday", "appointment_id", appointment.ID, "holiday", holiday.Name)
			return true, nil
		}
	}
	return false, nil
}

// retrieverPublicHolidays retrieves public holidays for a given year and country.
// It first checks the cache, and if not found, fetches from an external API
// and stores the result in the cache.
func (a AppointmentDefaultUseCase) retrieverPublicHolidays(ctx context.Context, year int, country string) ([]entity.PublicHolidays, error) {
	// 1. Check cache
	holidays, err := a.holidaysCacheClient.GetPublicHolidays(ctx, year, country)
	if err == nil {
		a.logger.Trace(ctx, "cache hit for public holidays", "year", year, "country", country)
		return holidays, nil
	}

	// 2. Fetch from external API
	holidays, err = a.holidaysHttpClient.RetrievePublicHolidays(year, country)
	if err != nil {
		a.logger.Error(ctx, err, "failed to retrieve public holidays from external API")
		return nil, apperr.BadGateway("failed to retrieve public holidays", err)
	}

	// 3 . Store in cache
	if err := a.holidaysCacheClient.SetPublicHolidays(ctx, year, country, holidays); err != nil {
		a.logger.Trace(ctx, "cache set failed for public holidays", "year", year, "country", country)
		a.logger.Error(ctx, err, "failed to store public holidays in cache")
	}
	return holidays, nil
}
