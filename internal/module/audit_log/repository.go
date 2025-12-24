package audit_log

import "gorm.io/gorm"

// Repository bertanggung jawab menyimpan audit log
// NOTE: append-only, tidak ada update / delete
type Repository interface {
	Create(log *AuditLog) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// ========================
// CREATE AUDIT LOG
// ========================
func (r *repository) Create(log *AuditLog) error {
	return r.db.Create(log).Error
}
