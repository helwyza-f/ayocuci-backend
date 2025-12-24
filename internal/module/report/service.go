package report

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// ========================
// RESPONSE DTO
// ========================
type Summary struct {
	Date         string `json:"date,omitempty"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
	TotalSales   int64  `json:"total_sales"`
	TotalExpense int64  `json:"total_expense"`
	TotalOrders  int64  `json:"total_orders"`
	Profit       int64  `json:"profit"`
}

type Service interface {
	Daily(outletID uint, date time.Time) (*Summary, error)
	Range(outletID uint, start time.Time, end time.Time) (*Summary, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db}
}

// ========================
// DAILY REPORT
// ========================
func (s *service) Daily(outletID uint, date time.Time) (*Summary, error) {
	start := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		date.Location(),
	)
	end := start.Add(24 * time.Hour)
	log.Print(start, end)
	return s.aggregate(outletID, start, end, true)
}

// ========================
// RANGE REPORT
// ========================
func (s *service) Range(outletID uint, start time.Time, end time.Time) (*Summary, error) {
	return s.aggregate(outletID, start, end, false)
}

// ========================
// AGGREGATION CORE
// ========================
func (s *service) aggregate(
	outletID uint,
	start time.Time,
	end time.Time,
	daily bool,
) (*Summary, error) {

	var totalSales int64
	var totalExpense int64
	var totalOrders int64

	// ========================
	// TOTAL SALES
	// ========================
	s.db.
		Table("orders").
		Where(
			"outlet_id = ? AND created_at BETWEEN ? AND ?",
			outletID,
			start,
			end,
		).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalSales)

	// ========================
	// TOTAL ORDERS
	// ========================
	s.db.
		Table("orders").
		Where(
			"outlet_id = ? AND created_at BETWEEN ? AND ?",
			outletID,
			start,
			end,
		).
		Count(&totalOrders)

	// ========================
	// TOTAL EXPENSE (DATETIME SAFE)
	// ========================
	s.db.
		Table("expenses").
		Where(
			"outlet_id = ? AND created_at BETWEEN ? AND ?",
			outletID,
			start,
			end,
		).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	// ========================
	// RESULT
	// ========================
	result := &Summary{
		TotalSales:   totalSales,
		TotalExpense: totalExpense,
		TotalOrders:  totalOrders,
		Profit:       totalSales - totalExpense,
	}

	if daily {
		result.Date = start.Format("2006-01-02")
	} else {
		result.StartDate = start.Format("2006-01-02")
		result.EndDate = end.Format("2006-01-02")
	}

	return result, nil
}
