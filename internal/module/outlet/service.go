package outlet

import "errors"

type Service interface {

	List(clientID uint, role string) ([]Outlet, error)
	Create(clientID uint, userID uint, name, address, phone string) (*Outlet, error)
	Delete(id uint, clientID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Delete(id uint, clientID uint) error {
	return s.repo.Delete(id, clientID)
}

func (s *service) Create(
	clientID uint,
	userID uint,
	name, address, phone string,
) (*Outlet, error) {

	if clientID == 0 {
		return nil, errors.New("invalid client")
	}
	if name == "" {
		return nil, errors.New("name required")
	}

	outlet := &Outlet{
		ClientID: clientID,
		Name:     name,
		Address:  address,
		Phone:    phone,
		Active:   true,
	}

	// Gunakan fungsi transaksi
    if err := s.repo.CreateWithOwner(outlet, userID); err != nil {
        return nil, err
    }

	return outlet, nil
}


func (s *service) List(clientID uint, role string) ([]Outlet, error) {
	if role == "super_admin" {
		return s.repo.FindAll()
	}
	return s.repo.FindByClient(clientID)
}
