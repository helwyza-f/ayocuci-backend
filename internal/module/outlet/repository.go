package outlet

import "gorm.io/gorm"

type Repository interface {
	Create(outlet *Outlet) error
	FindByClient(clientID uint) ([]Outlet, error)
	FindAll() ([]Outlet, error) // untuk super_admin
	FindByID(id uint) (*Outlet, error)
	CreateWithOwner(outlet *Outlet, userID uint) error
	Delete(id uint, clientID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(outlet *Outlet) error {
	return r.db.Create(outlet).Error
}

func (r *repository) FindByClient(clientID uint) ([]Outlet, error) {
    var list []Outlet
    // ⚠️ WAJIB: Filter Active = true
    err := r.db.Where("client_id = ? AND active = ?", clientID, true).Find(&list).Error
    return list, err
}

func (r *repository) Delete(id uint, clientID uint) error {
    // Kita pastikan client_id cocok supaya owner tidak bisa hapus outlet orang lain
    return r.db.Model(&Outlet{}).
        Where("id = ? AND client_id = ?", id, clientID).
        Update("active", false).Error
}

func (r *repository) FindAll() ([]Outlet, error) {
	var outlets []Outlet
	err := r.db.Find(&outlets).Error
	return outlets, err
}

func (r *repository) FindByID(id uint) (*Outlet, error) {
	var outlet Outlet
	err := r.db.First(&outlet, id).Error
	return &outlet, err
}

func (r *repository) CreateWithOwner(outlet *Outlet, userID uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. Simpan Outlet
        if err := tx.Create(outlet).Error; err != nil {
            return err
        }

        // 2. Assign Owner ke tabel user_outlets
        // Kita langsung tembak ke tabel user_outlets
        userOutlet := map[string]interface{}{
            "user_id":   userID,
            "outlet_id": outlet.ID,
            "role":      "owner",
            "active":    true,
        }

        if err := tx.Table("user_outlets").Create(&userOutlet).Error; err != nil {
            return err
        }

        return nil
    })
}
