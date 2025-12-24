package common

import "gorm.io/gorm"

func DataScope(role string, outletID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if role == "super_admin" {
			return db
		}
		return db.Where("outlet_id = ?", outletID)
	}
}
