package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
	"tabeo.org/challenge/internal/pkg/logger"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, a *entity.Appointment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

type mockCache struct{ mock.Mock }

func (m *mockCache) GetPublicHolidays(ctx context.Context, year int, country string) ([]entity.PublicHolidays, error) {
	args := m.Called(ctx, year, country)
	var res []entity.PublicHolidays
	if v := args.Get(0); v != nil {
		res = v.([]entity.PublicHolidays)
	}
	return res, args.Error(1)
}
func (m *mockCache) SetPublicHolidays(ctx context.Context, year int, country string, holidays []entity.PublicHolidays) error {
	args := m.Called(ctx, year, country, holidays)
	return args.Error(0)
}

type mockHolidayClient struct{ mock.Mock }

func (m *mockHolidayClient) RetrievePublicHolidays(year int, country string) ([]entity.PublicHolidays, error) {
	args := m.Called(year, country)
	var res []entity.PublicHolidays
	if v := args.Get(0); v != nil {
		res = v.([]entity.PublicHolidays)
	}
	return res, args.Error(1)
}

type mockLogger struct{ mock.Mock }

func (m *mockLogger) Error(ctx context.Context, err error, msg string, kv ...any) {}
func (m *mockLogger) Warn(ctx context.Context, msg string, kv ...any)             {}
func (m *mockLogger) Trace(ctx context.Context, msg string, kv ...any)            {}
func (m *mockLogger) Debug(ctx context.Context, msg string, kv ...any)            {}
func (m *mockLogger) Info(ctx context.Context, msg string, kv ...any)             {}
func (m *mockLogger) Fatal(ctx context.Context, err error, msg string, kv ...any) {}
func (m *mockLogger) With(kv ...any) logger.AppLogger                             { return m }

func makeUsecase(repo *mockRepo, cache *mockCache, client *mockHolidayClient, log *mockLogger) *AppointmentDefaultUseCase {
	return &AppointmentDefaultUseCase{
		appointmentRepo:     repo,
		holidaysCacheClient: cache,
		holidaysHttpClient:  client,
		logger:              log,
	}
}

func TestCreateAppointment_Success(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	appt := &entity.Appointment{ID: [16]byte{1}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{}, nil)
	repo.On("Create", mock.Anything, appt).Return(nil)

	id, err := uc.CreateAppointment(context.Background(), appt)
	assert.NoError(t, err)
	assert.Equal(t, appt.ID.String(), id)
}

func TestCreateAppointment_HolidayConflict(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	appt := &entity.Appointment{ID: [16]byte{2}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	holiday := entity.PublicHolidays{Date: "2025-01-01", Name: "New Year"}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{holiday}, nil)

	id, err := uc.CreateAppointment(context.Background(), appt)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.True(t, apperr.IsInvalid(err))
}

func TestCreateAppointment_RepoError(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	appt := &entity.Appointment{ID: [16]byte{3}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{}, nil)
	repo.On("Create", mock.Anything, appt).Return(errors.New("db error"))

	id, err := uc.CreateAppointment(context.Background(), appt)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.True(t, apperr.IsInternal(err))
}

func TestCreateAppointment_PublicHolidayError_ContinuesToCreate(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	appt := &entity.Appointment{ID: [16]byte{4}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return(nil, errors.New("cache error"))
	client.On("RetrievePublicHolidays", 2025, "GB").Return(nil, errors.New("api error"))
	repo.On("Create", mock.Anything, appt).Return(nil)

	id, err := uc.CreateAppointment(context.Background(), appt)
	assert.NoError(t, err)
	assert.Equal(t, appt.ID.String(), id)
}

func TestDismissAppointment_OnHoliday(t *testing.T) {
	viper.Set("holidays.country", "GB")
	cache := &mockCache{}
	uc := &AppointmentDefaultUseCase{holidaysCacheClient: cache, logger: &mockLogger{}}

	appt := &entity.Appointment{ID: [16]byte{5}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{{Date: "2025-01-01", Name: "New Year"}}, nil)

	result, err := uc.dismissAppointment(context.Background(), appt)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestDismissAppointment_NotOnHoliday(t *testing.T) {
	viper.Set("holidays.country", "GB")
	cache := &mockCache{}
	uc := &AppointmentDefaultUseCase{holidaysCacheClient: cache, logger: &mockLogger{}}

	appt := &entity.Appointment{ID: [16]byte{6}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{{Date: "2025-01-01", Name: "New Year"}}, nil)

	result, err := uc.dismissAppointment(context.Background(), appt)
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestDismissAppointment_ParseError(t *testing.T) {
	viper.Set("holidays.country", "GB")
	cache := &mockCache{}
	uc := &AppointmentDefaultUseCase{holidaysCacheClient: cache, logger: &mockLogger{}}

	appt := &entity.Appointment{ID: [16]byte{7}, FirstName: "A", LastName: "B", VisitDate: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)}
	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{{Date: "bad-date", Name: "Bad Date"}}, nil)

	result, err := uc.dismissAppointment(context.Background(), appt)
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestRetrieverPublicHolidays_CacheHit(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return([]entity.PublicHolidays{{Date: "2025-01-01"}}, nil)

	result, err := uc.retrieverPublicHolidays(context.Background(), 2025, "GB")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestRetrieverPublicHolidays_CacheMissAndApiSuccess(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return(nil, errors.New("cache miss"))
	client.On("RetrievePublicHolidays", 2025, "GB").Return([]entity.PublicHolidays{{Date: "2025-01-01"}}, nil)
	cache.On("SetPublicHolidays", mock.Anything, 2025, "GB", mock.Anything).Return(nil)

	result, err := uc.retrieverPublicHolidays(context.Background(), 2025, "GB")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestRetrieverPublicHolidays_CacheMissAndApiError(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return(nil, errors.New("cache miss"))
	client.On("RetrievePublicHolidays", 2025, "GB").Return(nil, errors.New("api error"))

	result, err := uc.retrieverPublicHolidays(context.Background(), 2025, "GB")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apperr.IsBadGateway(err))
}

func TestRetrieverPublicHolidays_CacheSetError(t *testing.T) {
	viper.Set("holidays.country", "GB")
	repo := &mockRepo{}
	cache := &mockCache{}
	client := &mockHolidayClient{}
	log := &mockLogger{}
	uc := makeUsecase(repo, cache, client, log)

	cache.On("GetPublicHolidays", mock.Anything, 2025, "GB").Return(nil, errors.New("cache miss"))
	client.On("RetrievePublicHolidays", 2025, "GB").Return([]entity.PublicHolidays{{Date: "2025-01-01"}}, nil)
	cache.On("SetPublicHolidays", mock.Anything, 2025, "GB", mock.Anything).Return(errors.New("cache set error"))

	result, err := uc.retrieverPublicHolidays(context.Background(), 2025, "GB")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
