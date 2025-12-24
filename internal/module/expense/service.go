package expense

import (
	"errors"
	"time"
)

type Service interface {
	Create(outletID uint, category string, amount int64, note string, date time.Time) error
	GetAll(outletID uint) ([]Expense, error)
	Update(outletID uint, id uint, category string, amount int64, note string, date time.Time) error
	Delete(outletID uint, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// ========================
// CREATE EXPENSE
// ========================
func (s *service) Create(
	outletID uint,
	category string,
	amount int64,
	note string,
	date time.Time,
) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	if outletID == 0 {
		return errors.New("outlet context required")
	}
	if category == "" {
		return errors.New("category required")
	}
	if amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if date.IsZero() {
		date = time.Now().In(loc)
	}

	expense := &Expense{
		OutletID: outletID,
		Category: category,
		Amount:   amount,
		Note:     note,
		Date:     date,
	}

	return s.repo.Create(expense)
}

// ========================
// LIST EXPENSES
// ========================
func (s *service) GetAll(outletID uint) ([]Expense, error) {
	if outletID == 0 {
		return nil, errors.New("outlet context required")
	}
	return s.repo.FindAll(outletID)
}

// ========================
// UPDATE EXPENSE
// ========================
func (s *service) Update(
	outletID uint,
	id uint,
	category string,
	amount int64,
	note string,
	date time.Time,
) error {

	if outletID == 0 {
		return errors.New("outlet context required")
	}

	expense, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return err
	}

	if category != "" {
		expense.Category = category
	}
	if amount > 0 {
		expense.Amount = amount
	}
	expense.Note = note
	if !date.IsZero() {
		expense.Date = date
	}

	return s.repo.Update(expense)
}

// ========================
// DELETE EXPENSE
// ========================
func (s *service) Delete(outletID uint, id uint) error {
	if outletID == 0 {
		return errors.New("outlet context required")
	}
	return s.repo.Delete(outletID, id)
}
