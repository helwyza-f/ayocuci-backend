package laundry_service

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	OutletID uint   `json:"outlet_id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Estimate string `json:"estimate"` // contoh: "2 hari"
	Active   bool   `json:"active"`
}
