package customer

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	OutletID uint   `json:"outlet_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}
