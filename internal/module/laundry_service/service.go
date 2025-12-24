package laundry_service

import "errors"

type ServiceLogic interface {
	Create(outletID uint, name string, price int64, estimate string) error
	GetAll(outletID uint) ([]Service, error)
	Update(outletID uint, id uint, name string, price int64, estimate string, active bool) error
	Delete(outletID uint, id uint) error
}

type serviceLogic struct {
	repo Repository
}

func NewService(repo Repository) ServiceLogic {
	return &serviceLogic{repo}
}

func (s *serviceLogic) Create(outletID uint, name string, price int64, estimate string) error {
	if outletID == 0 {
		return errors.New("outlet context required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}

	service := &Service{
		OutletID: outletID,
		Name:     name,
		Price:    price,
		Estimate: estimate,
		Active:   true,
	}

	return s.repo.Create(service)
}

func (s *serviceLogic) GetAll(outletID uint) ([]Service, error) {
	if outletID == 0 {
		return nil, errors.New("outlet context required")
	}
	return s.repo.FindAll(outletID)
}

func (s *serviceLogic) Update(
	outletID uint,
	id uint,
	name string,
	price int64,
	estimate string,
	active bool ,
) error {

	if outletID == 0 {
		return errors.New("outlet context required")
	}

	service, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return err
	}

	if name != "" {
		service.Name = name
	}
	if price > 0 {
		service.Price = price
	}

	service.Estimate = estimate
	service.Active = active

	return s.repo.Update(service)
}

func (s *serviceLogic) Delete(
	outletID uint,
	id uint,
) error {

	if outletID == 0 {
		return errors.New("outlet context required")
	}
	service, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return err
	}
	service.Active = false

	return s.repo.Update(service)
}
