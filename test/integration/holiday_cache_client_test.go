package integration

import (
	"context"
	"github.com/spf13/viper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"tabeo.org/challenge/internal/adapter/cache"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
)

func TestHolidayCacheClient_SetAndGetPublicHolidays(t *testing.T) {
	initConfig()
	viper.Set("cache.db", 0)
	viper.Set("cache.ttl", 2)

	client := cache.NewHolidayCacheClient()
	ctx := context.Background()
	year := 2025
	country := "DE"
	holiday := entity.PublicHolidays{
		Date:        "2025-01-01",
		LocalName:   "Neujahr",
		Name:        "New Year's Day",
		CountryCode: country,
		Fixed:       true,
		Global:      true,
		Counties:    nil,
		LaunchYear:  1967,
		Types:       []string{"Public"},
	}
	holidays := []entity.PublicHolidays{holiday}

	_, err := client.GetPublicHolidays(ctx, year, country)
	assert.Error(t, err)
	assert.True(t, apperr.IsNotFound(err))

	err = client.SetPublicHolidays(ctx, year, country, holidays)
	assert.NoError(t, err)

	got, err := client.GetPublicHolidays(ctx, year, country)
	assert.NoError(t, err)
	assert.Equal(t, holidays, got)

	time.Sleep(3 * time.Second)
	_, err = client.GetPublicHolidays(ctx, year, country)
	assert.Error(t, err)
	assert.True(t, apperr.IsNotFound(err))
}

func TestHolidayCacheClient_SetInvalidData(t *testing.T) {
	initConfig()
	viper.Set("cache.db", 0)

	client := cache.NewHolidayCacheClient()
	ctx := context.Background()
	year := 2025
	country := "FR"

	err := client.SetPublicHolidays(ctx, year, country, nil)
	assert.NoError(t, err)

	got, err := client.GetPublicHolidays(ctx, year, country)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestHolidayCacheClient_Integration_ConnectToExistingCache(t *testing.T) {
	initConfig()
	viper.Set("cache.db", 0)

	client := cache.NewHolidayCacheClient()
	ctx := context.Background()
	year := 2026
	country := "IT"

	_, err := client.GetPublicHolidays(ctx, year, country)
	assert.Error(t, err)
	assert.True(t, apperr.IsNotFound(err))

	holiday := entity.PublicHolidays{
		Date:        "2026-04-25",
		LocalName:   "Festa della Liberazione",
		Name:        "Liberation Day",
		CountryCode: country,
		Fixed:       true,
		Global:      false,
		Counties:    nil,
		LaunchYear:  1946,
		Types:       []string{"Public"},
	}
	holidays := []entity.PublicHolidays{holiday}

	err = client.SetPublicHolidays(ctx, year, country, holidays)
	assert.NoError(t, err)

	got, err := client.GetPublicHolidays(ctx, year, country)
	assert.NoError(t, err)
	assert.Equal(t, holidays, got)
}
