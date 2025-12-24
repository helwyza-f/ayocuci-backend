package expense

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	OutletID uint      `json:"outlet_id"`
	Category string    `json:"category"`
	Amount   int64     `json:"amount"`
	Note     string    `json:"note"`
	Date     time.Time `json:"date"`
}
