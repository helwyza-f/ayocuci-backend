package audit_log

import "gorm.io/gorm"

// AuditLog menyimpan histori perubahan data (append-only)
type AuditLog struct {
	gorm.Model

	// multi-tenant context
	ClientID uint `json:"client_id" gorm:"index;not null"`
	OutletID uint `json:"outlet_id" gorm:"index"`
	UserID   uint `json:"user_id" gorm:"index"`

	// apa yang terjadi
	Action string `json:"action" gorm:"type:varchar(50);not null"`
	Entity string `json:"entity" gorm:"type:varchar(50);not null"`
	EntityID uint `json:"entity_id" gorm:"index"`

	// snapshot data
	OldData string `json:"old_data" gorm:"type:json"`
	NewData string `json:"new_data" gorm:"type:json"`
}
