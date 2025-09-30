package repo

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

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

// Create inserts a new appointment into the database
func (r *AppointmentDefaultRepository) Create(ctx context.Context, a *entity.Appointment) error {
	a.ID = uuid.New()
	if err := r.DB.WithContext(ctx).Create(a).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			return apperr.Exists("visit date already booked", err)
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return apperr.Internal("create appointment failed", err)
	}
	return nil
}
