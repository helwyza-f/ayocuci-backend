package expense

import "gorm.io/gorm"

type Repository interface {
	Create(expense *Expense) error
	FindAll(outletID uint) ([]Expense, error)
	FindByID(outletID uint, id uint) (*Expense, error)
	Update(expense *Expense) error
	Delete(outletID uint, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(expense *Expense) error {
	return r.db.Create(expense).Error
}

func (r *repository) FindAll(outletID uint) ([]Expense, error) {
	var expenses []Expense
	err := r.db.
		Where("outlet_id = ?", outletID).
		Order("date desc").
		Find(&expenses).Error
	return expenses, err
}

func (r *repository) FindByID(outletID uint, id uint) (*Expense, error) {
	var expense Expense
	err := r.db.
		Where("outlet_id = ?", outletID).
		First(&expense, id).Error
	return &expense, err
}

func (r *repository) Update(expense *Expense) error {
	return r.db.Save(expense).Error
}

func (r *repository) Delete(outletID uint, id uint) error {
	return r.db.
		Where("outlet_id = ?", outletID).
		Delete(&Expense{}, id).Error
}
