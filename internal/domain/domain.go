package domain

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"column:user_id; type:uuid"`
	ServiceName string    `json:"service_name" gorm:"service_name"`
	Price       int       `json:"price" gorm:"price"`
	StartDate   time.Time `json:"start_date" gorm:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty" gorm:"end_date"`
}
