package entity

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID        uuid.UUID `gorm:"type:binary(16);primaryKey"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	VisitDate time.Time `gorm:"type:date;unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
