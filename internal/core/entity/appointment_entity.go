package entity

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	VisitDate time.Time `gorm:"type:date;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
}

func (Appointment) TableName() string {
	return "tabeo.appointment"
}
