package client

import "time"

type Client struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
}
