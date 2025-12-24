package order

import (
	"errors"

	"github.com/helwyza-f/ayocuci-backend/internal/module/audit_log"
)

type ItemInput struct {
	ServiceID uint    `json:"service_id"`
	Qty       float64 `json:"qty"`
	Price     int64   `json:"price"`
}

type Service interface {
	Create(outletID, customerID uint, items []ItemInput) error
	GetAll(outletID uint, filter OrderFilter) ([]Order, int64, error)
	UpdateStatus(
		clientID uint,
		outletID uint,
		userID uint,
		orderID uint,
		status string,
	) error
}

type service struct {
	repo        Repository
	auditLogger audit_log.Service
}

func NewService(
	repo Repository,
	auditLogger audit_log.Service,
) Service {
	return &service{
		repo:        repo,
		auditLogger: auditLogger,
	}
}

// ========================
// CREATE ORDER
// (audit bisa ditambah nanti)
// ========================
func (s *service) Create(outletID, customerID uint, items []ItemInput) error {
	if outletID == 0 {
		return errors.New("outlet context required")
	}
	if customerID == 0 {
		return errors.New("customer_id required")
	}
	if len(items) == 0 {
		return errors.New("items required")
	}

	var total int64
	var orderItems []OrderItem

	for _, item := range items {
		if item.ServiceID == 0 {
			return errors.New("service_id required")
		}
		if item.Qty <= 0 || item.Price <= 0 {
			return errors.New("invalid qty or price")
		}

		total += int64(item.Qty * float64(item.Price))
		orderItems = append(orderItems, OrderItem{
			ServiceID: item.ServiceID,
			Qty:       item.Qty,
			Price:     item.Price,
		})
	}

	order := &Order{
		OutletID:   outletID,
		CustomerID: customerID,
		Status:     "baru",
		Total:      total,
	}

	return s.repo.Create(order, orderItems)
}

// ========================
// LIST ORDERS
// ========================
func (s *service) GetAll(
	outletID uint,
	filter OrderFilter,
) ([]Order, int64, error) {

	if outletID == 0 {
		return nil, 0, errors.New("outlet context required")
	}

	return s.repo.FindAll(outletID, filter)
}

// ========================
// UPDATE STATUS + AUDIT LOG
// ========================
func (s *service) UpdateStatus(
	clientID uint,
	outletID uint,
	userID uint,
	orderID uint,
	status string,
) error {

	if outletID == 0 {
		return errors.New("outlet context required")
	}

	valid := map[string]bool{
		"baru":     true,
		"diproses": true,
		"selesai":  true,
		"diambil":  true,
	}
	if !valid[status] {
		return errors.New("invalid status")
	}

	// ambil data lama
	order, err := s.repo.FindByID(outletID, orderID)
	if err != nil {
		return err
	}

	oldOrder := *order // snapshot sebelum berubah

	// update
	order.Status = status
	if err := s.repo.UpdateStatus(order); err != nil {
		return err
	}

	// 🔥 AUDIT LOG (NON-BLOCKING)
	_ = s.auditLogger.Log(
		clientID,
		outletID,
		userID,
		"STATUS_CHANGE",
		"order",
		order.ID,
		oldOrder,
		order,
	)

	return nil
}
