package customer

import "errors"

type Service interface {
	Create(outletID uint, name, phone string) error
	GetAll(outletID uint) ([]Customer, error)
	GetByID(outletID uint, id uint) (*Customer, error)
	Update(outletID uint, id uint, name, phone string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(outletID uint, name, phone string) error {
	if outletID == 0 {
		return errors.New("outlet context required")
	}
	if name == "" {
		return errors.New("name required")
	}

	customer := &Customer{
		OutletID: outletID,
		Name:     name,
		Phone:    phone,
	}

	return s.repo.Create(customer)
}

func (s *service) GetAll(outletID uint) ([]Customer, error) {
	if outletID == 0 {
		return nil, errors.New("outlet context required")
	}
	return s.repo.FindAll(outletID)
}

func (s *service) GetByID(outletID uint, id uint) (*Customer, error) {
	if outletID == 0 {
		return nil, errors.New("outlet context required")
	}
	return s.repo.FindByID(outletID, id)
}

func (s *service) Update(outletID uint, id uint, name, phone string) error {
	customer, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return err
	}

	if name != "" {
		customer.Name = name
	}
	customer.Phone = phone

	return s.repo.Update(customer)
}
