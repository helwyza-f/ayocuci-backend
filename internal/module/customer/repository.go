package customer

import "gorm.io/gorm"

type Repository interface {
	Create(customer *Customer) error
	FindAll(outletID uint) ([]Customer, error)
	FindByID(outletID uint, id uint) (*Customer, error)
	Update(customer *Customer) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *repository) FindAll(outletID uint) ([]Customer, error) {
	var customers []Customer
	err := r.db.
		Where("outlet_id = ?", outletID).
		Order("created_at desc").
		Find(&customers).Error
	return customers, err
}

func (r *repository) FindByID(outletID uint, id uint) (*Customer, error) {
	var customer Customer
	err := r.db.
		Where("outlet_id = ?", outletID).
		First(&customer, id).Error
	return &customer, err
}

func (r *repository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}
