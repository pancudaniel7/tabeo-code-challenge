package cache

import (
	"context"
	"tabeo.org/challenge/internal/core/entity"
)

// HolidayCacheClient defines the interface for caching public holidays
// by year and country.
type HolidayCacheClient interface {
	GetPublicHolidays(ctx context.Context, year int, country string) ([]entity.PublicHolidays, error)
	SetPublicHolidays(ctx context.Context, year int, country string, holidays []entity.PublicHolidays) error
}
