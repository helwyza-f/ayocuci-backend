package auth

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	UpdateActiveOutlet(userID, outletID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	
	return &user, err
}

func (r *repository) UpdateActiveOutlet(userID, outletID uint) error {
	return r.db.Model(&User{}).
		Where("id = ?", userID).
		Update("active_outlet_id", outletID).
		Error
}

