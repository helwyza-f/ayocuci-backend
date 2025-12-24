package order

import "gorm.io/gorm"

// =======================
// ORDER (HEADER)
// =======================
type Order struct {
	gorm.Model

	OutletID   uint   `json:"outlet_id" gorm:"index;not null"`
	CustomerID uint   `json:"customer_id" gorm:"index;not null"`

	Status string `json:"status" gorm:"type:varchar(20);not null"`
	Total  int64  `json:"total" gorm:"not null"`

	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// =======================
// ORDER ITEM (DETAIL)
// =======================
type OrderItem struct {
	gorm.Model

	OrderID   uint    `json:"order_id" gorm:"index;not null"`
	ServiceID uint    `json:"service_id" gorm:"index;not null"`

	Qty   float64 `json:"qty" gorm:"not null"`
	Price int64   `json:"price" gorm:"not null"`
}
