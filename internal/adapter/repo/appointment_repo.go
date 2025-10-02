package repo

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
)

// AppointmentDefaultRepository provides methods to interact with appointments in the DB
type AppointmentDefaultRepository struct {
	DB *gorm.DB
}

func NewAppointmentDefaultRepository(db *gorm.DB) *AppointmentDefaultRepository {
	return &AppointmentDefaultRepository{DB: db}
}

func (r *AppointmentDefaultRepository) FindByVistDate(ctx context.Context, appointmentDate time.Time) (*entity.Appointment, error) {
	var appointment entity.Appointment
	if err := r.DB.WithContext(ctx).Where("visit_date = ?", appointmentDate).First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFoundErr("appointment not found", err)
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, apperr.Internal("operation canceled", err)
		}
		return nil, apperr.Internal("find appointment by visit date failed", err)
	}
	return &appointment, nil
}

// Create inserts a new appointment into the database
func (r *AppointmentDefaultRepository) Create(ctx context.Context, a *entity.Appointment) error {
	a.ID = uuid.New()
	if err := r.DB.WithContext(ctx).Create(a).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			return apperr.Exists("visit date already booked", err)
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return apperr.Internal("operation canceled", err)
		}
		return apperr.Internal("create appointment failed", err)
	}
	return nil
}
