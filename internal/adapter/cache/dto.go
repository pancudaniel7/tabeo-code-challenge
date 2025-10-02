package cache

import "tabeo.org/challenge/internal/core/entity"

type PublicHolidaysCacheDTO struct {
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

func (dto *PublicHolidaysCacheDTO) ToEntity() entity.PublicHolidays {
	return entity.PublicHolidays{
		Date:        dto.Date,
		LocalName:   dto.LocalName,
		Name:        dto.Name,
		CountryCode: dto.CountryCode,
		Fixed:       dto.Fixed,
		Global:      dto.Global,
		Counties:    dto.Counties,
		LaunchYear:  dto.LaunchYear,
		Types:       dto.Types,
	}
}

func (dto *PublicHolidaysCacheDTO) ToDTO(entity *entity.PublicHolidays) *PublicHolidaysCacheDTO {
	return &PublicHolidaysCacheDTO{
		Date:        entity.Date,
		LocalName:   entity.LocalName,
		Name:        entity.Name,
		CountryCode: entity.CountryCode,
		Fixed:       entity.Fixed,
		Global:      entity.Global,
		Counties:    entity.Counties,
		LaunchYear:  entity.LaunchYear,
		Types:       entity.Types,
	}
}
