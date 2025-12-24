package employee

import (
	"errors"

	"github.com/helwyza-f/ayocuci-backend/internal/module/auth"

	"github.com/helwyza-f/ayocuci-backend/internal/module/outlet"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// buat akun employee (user login)
	CreateUser(clientID uint, email string, password string) (uint, error)

	// assign employee ke outlet
	Assign(requesterClientID uint, userID uint, outletID uint, role string) error

	// 🔥 ambil outlet yang bisa dipilih user
	GetMyOutlets(userID uint) ([]outlet.Outlet, error)

	MyOutlets(role string, userID uint, clientID uint) ([]outlet.Outlet, error)

	ListEmployees(outletID uint) ([]Employee, error)
	Transfer(clientID, userID, targetOutletID uint, role string) error
}

type service struct {
	repo       Repository        // user_outlets
	userRepo   auth.Repository   // users
	outletRepo outlet.Repository // outlets
	
}

func NewService(
	repo Repository,
	userRepo auth.Repository,
	outletRepo outlet.Repository,
	
) Service {
	return &service{
		repo:       repo,
		userRepo:   userRepo,
		outletRepo: outletRepo,
	
	}
}

// ========================
// CREATE EMPLOYEE USER
// ========================
func (s *service) CreateUser(
	clientID uint,
	email string,
	password string,
) (uint, error) {

	if email == "" || password == "" {
		return 0, errors.New("invalid input")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user := &auth.User{
		Email:    email,
		Password: string(hashed),
		Role:     "staff",
		ClientID: clientID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

// ========================
// ASSIGN EMPLOYEE TO OUTLET
// ========================
func (s *service) Transfer(clientID, userID, targetOutletID uint, role string) error {
    // 1. Validasi tenant (sama kayak assign, pastiin outlet tujuan milik si owner)
    outlet, err := s.outletRepo.FindByID(targetOutletID)
    if err != nil || outlet.ClientID != clientID {
        return errors.New("invalid target outlet")
    }

    // 2. Eksekusi Update di Repo
    return s.repo.UpdatePlacement(userID, targetOutletID, role)
}

func (s *service) Assign(
	requesterClientID uint,
	userID uint,
	outletID uint,
	role string,
) error {

	// 1️⃣ ambil user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Role != "staff" {
		return errors.New("only staff users can be assigned as employees")
	}

	// 2️⃣ ambil outlet
	outletData, err := s.outletRepo.FindByID(outletID)
	if err != nil {
		return errors.New("outlet not found")
	}

	// 3️⃣ validasi multi-tenant
	if user.ClientID != requesterClientID {
		return errors.New("cannot assign user from another client")
	}
	if outletData.ClientID != requesterClientID {
		return errors.New("cannot assign to outlet from another client")
	}

	// 4️⃣ 🔥 CEK DUPLIKAT
	exists, err := s.repo.Exists(userID, outletID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already assigned to this outlet")
	}

	// 5️⃣ assign
	return s.repo.CreateUser(userID, outletID, role)
}

// Di internal/module/employee/service.go
func (s *service) ListEmployees(outletID uint) ([]Employee, error) {
    if outletID == 0 {
        return nil, errors.New("outlet context required")
    }
    return s.repo.FindByOutlet(outletID)
}

// ========================
// LIST OUTLET USER (PILIH OUTLET AKTIF)
// ========================
func (s *service) GetMyOutlets(userID uint) ([]outlet.Outlet, error) {
	return s.repo.FindOutletsByUser(userID)
}



func (s *service) MyOutlets(
	role string,
	userID uint,
	clientID uint,
) ([]outlet.Outlet, error) {

	if role == "owner" {
		return s.outletRepo.FindByClient(clientID)
	}

	// staff / kasir / kurir
	return s.repo.FindOutletsByUser(userID)
}

