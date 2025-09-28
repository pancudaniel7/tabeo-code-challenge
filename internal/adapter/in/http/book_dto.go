package http

import "tabeo.org/challenge/internal/core/entity"

type BookRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	VisitDate string `json:"visitDate"`
}

type BookResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	VisitDate string `json:"visitDate"`
}

func (r *BookRequest) ToEntity() *entity.Book {
	return &entity.Book{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		VisitDate: r.VisitDate,
	}
}
