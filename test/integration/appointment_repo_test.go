package integration

import (
	"context"
	"tabeo.org/challenge/internal/adapter/repo"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func withTx(t *testing.T) *gorm.DB {
	tx := testDB.Begin()
	t.Cleanup(func() { _ = tx.Rollback() })
	return tx
}

func TestAppointmentRepository_Create_Success_Integration(t *testing.T) {
	tx := withTx(t)
	r := repo.NewAppointmentDefaultRepository(tx)
	a := &entity.Appointment{
		FirstName: "Alice",
		LastName:  "Smith",
		VisitDate: time.Date(2075, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	err := r.Create(context.Background(), a)
	require.NoError(t, err)

	var count int64
	require.NoError(t, tx.Model(&entity.Appointment{}).Where("visit_date = ?", a.VisitDate).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func TestAppointmentRepository_Create_Duplicate_Integration(t *testing.T) {
	tx := withTx(t)
	r := repo.NewAppointmentDefaultRepository(tx)
	ctx := context.Background()
	d := time.Date(2075, 1, 3, 0, 0, 0, 0, time.UTC)

	require.NoError(t, r.Create(ctx, &entity.Appointment{FirstName: "Bob", LastName: "Jones", VisitDate: d}))
	err := r.Create(ctx, &entity.Appointment{FirstName: "Eve", LastName: "Miller", VisitDate: d})
	require.Error(t, err)
	require.True(t, apperr.IsExists(err))
}

func TestAppointmentRepository_Create_ContextCanceled_Integration(t *testing.T) {
	tx := withTx(t)
	r := repo.NewAppointmentDefaultRepository(tx)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	a := &entity.Appointment{
		FirstName: "Zoe",
		LastName:  "Lee",
		VisitDate: time.Date(2075, 1, 4, 0, 0, 0, 0, time.UTC),
	}
	err := r.Create(ctx, a)
	require.Error(t, err)
}
