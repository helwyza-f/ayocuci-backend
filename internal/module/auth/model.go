// internal/module/auth/model.go
package auth

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Email          string    `gorm:"uniqueIndex"`
	Password       string
	Role           string    // super_admin | owner | staff
	ClientID       uint
	ActiveOutletID *uint     `gorm:"index"` // ⬅️ nullable
	CreatedAt      time.Time
}

