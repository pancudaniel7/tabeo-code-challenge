package http

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"tabeo.org/challenge/internal/core/entity"
	"tabeo.org/challenge/internal/pkg/apperr"
)

type HolidayDefaultClient struct {
	client *http.Client
}

func NewHolidayClient() HolidayClient {
	return &HolidayDefaultClient{
		client: &http.Client{},
	}
}

func (h *HolidayDefaultClient) RetrievePublicHolidays(year int, country string) ([]entity.PublicHolidays, error) {
	url := viper.GetString("holidays.url")
	if url == "" {
		return nil, apperr.Internal("holiday api url is not set", nil)
	}

	endpoint := fmt.Sprintf(url, year, country)
	resp, err := h.client.Get(endpoint)
	if err != nil {
		return nil, apperr.Internal("failed to call holiday API", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperr.BadGateway(fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperr.Internal("failed to read response body", err)
	}

	var holidaysResp []PublicHolidaysResponse
	err = json.Unmarshal(body, &holidaysResp)
	if err != nil {
		return nil, apperr.Internal("failed to unmarshal holidays response", err)
	}

	holidays := make([]entity.PublicHolidays, len(holidaysResp))
	for i, h := range holidaysResp {
		holidays[i] = h.ToEntity()
	}
	return holidays, nil
}
