package outlet

import "time"

type Outlet struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ClientID  uint      `json:"client_id" gorm:"index"` // ⬅️ INI WAJIB
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
