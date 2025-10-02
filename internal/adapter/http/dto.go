package http

import (
	"time"

	"tabeo.org/challenge/internal/core/entity"
)

type AppointmentRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	VisitDate string `json:"visitDate"`
}

type AppointmentResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	VisitDate string `json:"visitDate"`
}

func (r *AppointmentRequest) ToEntity() (*entity.Appointment, error) {
	visitDate, err := time.Parse("2006-01-02", r.VisitDate)
	if err != nil {
		return nil, err
	}
	return &entity.Appointment{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		VisitDate: visitDate,
	}, nil
}

type PublicHolidaysResponse struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Fixed       bool     `json:"fixed"`
	Global      bool     `json:"global"`
	Counties    []string `json:"counties"`
	LaunchYear  int      `json:"launchYear"`
	Types       []string `json:"types"`
}
