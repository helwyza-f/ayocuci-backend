package auth

import (
	"errors"

	"github.com/helwyza-f/ayocuci-backend/internal/module/client"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	RegisterOwner(businessName, email, password string) error
	Login(email, password string) (*User, error)
}

type service struct {
	db         *gorm.DB
	userRepo   Repository
	clientRepo client.Repository
}

func NewService(
	db *gorm.DB,
	userRepo Repository,
	clientRepo client.Repository,
) Service {
	return &service{
		db:         db,
		userRepo:   userRepo,
		clientRepo: clientRepo,
	}
}

func (s *service) RegisterOwner(businessName, email, password string) error {
	if businessName == "" || email == "" || password == "" {
		return errors.New("invalid input")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. create client
		newClient := &client.Client{
			Name: businessName,
		}
		if err := tx.Create(newClient).Error; err != nil {
			return err
		}

		// 2. hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

		// 3. create owner user
		user := &User{
			Email:    email,
			Password: string(hashed),
			Role:     "owner",
			ClientID: newClient.ID,
		}

		return tx.Create(user).Error
	})
}

func (s *service) Login(email, password string) (*User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

