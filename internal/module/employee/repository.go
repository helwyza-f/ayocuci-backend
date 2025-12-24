package employee

import (
	"github.com/helwyza-f/ayocuci-backend/internal/module/outlet"
	"gorm.io/gorm"
)

type Repository interface {
	// assign user ke outlet
	CreateUser(userID, outletID uint, role string) error

	// list user_outlet (internal)
	FindByUser(userID uint) ([]UserOutlet, error)

	// cek duplikat assign
	Exists(userID, outletID uint) (bool, error)

	// 🔥 INI YANG KURANG
	// list outlet yang bisa diakses user
	FindOutletsByUser(userID uint) ([]outlet.Outlet, error)

	// list employee di outlet
	FindByOutlet(outletID uint) ([]Employee, error)

	// update penempatan outlet & role
	UpdatePlacement(userID, outletID uint, role string) error

}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateUser(userID, outletID uint, role string) error {
	emp := &UserOutlet{
		UserID:   userID,
		OutletID: outletID,
		Role:     role,
		Active:   true,
	}
	return r.db.Create(emp).Error
}

func (r *repository) FindByUser(userID uint) ([]UserOutlet, error) {
	var list []UserOutlet
	err := r.db.
		Where("user_id = ? AND active = true", userID).
		Find(&list).Error
	return list, err
}

func (r *repository) Exists(userID, outletID uint) (bool, error) {
	var count int64
	err := r.db.Model(&UserOutlet{}).
		Where("user_id = ? AND outlet_id = ? AND active = true", userID, outletID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) FindByOutlet(outletID uint) ([]Employee, error) {
    var list []Employee

    err := r.db.Table("user_outlets").
        Select("users.id, users.email, user_outlets.role, user_outlets.active").
        Joins("JOIN users ON users.id = user_outlets.user_id").
        // 🔥 Tambahkan filter AND user_outlets.role != 'owner'
        Where("user_outlets.outlet_id = ? AND user_outlets.active = ? AND user_outlets.role != ?", outletID, true, "owner").
        Scan(&list).Error

    return list, err
}

func (r *repository) UpdatePlacement(userID, outletID uint, role string) error {
    // Kita pakai Transaction supaya kalau gagal bisa rollback
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. Hapus/Non-aktifkan semua penempatan lama user ini
        // Agar tidak bentrok saat kita update ke outlet baru
        if err := tx.Where("user_id = ?", userID).Delete(&UserOutlet{}).Error; err != nil {
            return err
        }

        // 2. Buat penempatan baru
        newPlacement := UserOutlet{
            UserID:   userID,
            OutletID: outletID,
            Role:     role,
            Active:   true,
        }
        
        return tx.Create(&newPlacement).Error
    })
}

// 🔥 JOIN user_outlets → outlets
func (r *repository) FindOutletsByUser(userID uint) ([]outlet.Outlet, error) {
	var outlets []outlet.Outlet

	err := r.db.
		Table("outlets").
		Joins("JOIN user_outlets uo ON uo.outlet_id = outlets.id").
		Where("uo.user_id = ? AND uo.active = true", userID).
		Find(&outlets).Error

	return outlets, err
}
