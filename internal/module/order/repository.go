package order

import (
	"time"

	"gorm.io/gorm"
)

type OrderFilter struct {
	Status string
	From   *time.Time
	To     *time.Time
	Page   int
	Limit  int
}

type Repository interface {
	Create(order *Order, items []OrderItem) error
	FindAll(outletID uint, filter OrderFilter) ([]Order, int64, error)
	FindByID(outletID uint, id uint) (*Order, error)
	UpdateStatus(order *Order) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// ========================
// CREATE ORDER + ITEMS
// ========================
func (r *repository) Create(order *Order, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ========================
// LIST ORDERS (PAGINATED + FILTERED)
// ========================
func (r *repository) FindAll(
	outletID uint,
	filter OrderFilter,
) ([]Order, int64, error) {

	var (
		orders []Order
		total  int64
	)

	// base query (SELALU outlet scoped)
	query := r.db.
		Model(&Order{}).
		Where("outlet_id = ?", outletID)

	// filter status
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// filter tanggal
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}

	// count total sebelum pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// default pagination
	page := filter.Page
	limit := filter.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// ambil data
	err := query.
		Preload("Items").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

// ========================
// GET ORDER DETAIL
// ========================
func (r *repository) FindByID(outletID uint, id uint) (*Order, error) {
	var order Order

	err := r.db.
		Where("outlet_id = ?", outletID).
		Preload("Items").
		First(&order, id).Error

	return &order, err
}

// ========================
// UPDATE STATUS
// ========================
func (r *repository) UpdateStatus(order *Order) error {
	return r.db.Save(order).Error
}
