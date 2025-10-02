package integration

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"tabeo.org/challenge/internal/adapter/http"
	"tabeo.org/challenge/internal/pkg/apperr"
)

func TestRetrievePublicHolidays_Success(t *testing.T) {
	viper.Set("holiday.url", "https://date.nager.at/api/v3")
	client := http.NewHolidayClient()
	year := 2025
	country := "DE"

	holidays, err := client.RetrievePublicHolidays(year, country)
	assert.NoError(t, err)
	assert.NotEmpty(t, holidays)
	for _, h := range holidays {
		assert.Equal(t, country, h.CountryCode)
		assert.NotEmpty(t, h.Date)
		assert.NotEmpty(t, h.Name)
	}
}

func TestRetrievePublicHolidays_BadGateway(t *testing.T) {
	viper.Set("holiday.url", "https://date.nager.at/api/v3")
	client := http.NewHolidayClient()
	year := 2025
	country := "ZZZ"

	_, err := client.RetrievePublicHolidays(year, country)
	assert.Error(t, err)
	assert.True(t, apperr.IsBadGateway(err))
}

func TestRetrievePublicHolidays_Internal_NoURL(t *testing.T) {
	viper.Set("holiday.url", "")
	client := http.NewHolidayClient()
	year := 2025
	country := "DE"

	_, err := client.RetrievePublicHolidays(year, country)
	assert.Error(t, err)
	assert.True(t, apperr.IsInternal(err))
}
