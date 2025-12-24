package laundry_service

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(service *Service) error
	FindAll(outletID uint) ([]Service, error)
	FindByID(outletID uint, id uint) (*Service, error)
	Update(service *Service) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(service *Service) error {
	return r.db.Create(service).Error
}

func (r *repository) FindAll(outletID uint) ([]Service, error) {
	var services []Service
	log.Print(outletID)
	err := r.db.
		Where("outlet_id = ? AND active = ?", outletID, true).
		Order("created_at desc").
		Find(&services).Error
	return services, err
}

func (r *repository) FindByID(outletID uint, id uint) (*Service, error) {
	var service Service
	err := r.db.
		Where("outlet_id = ?", outletID).
		First(&service, id).Error
	return &service, err
}

func (r *repository) Update(service *Service) error {
	return r.db.Save(service).Error
}
