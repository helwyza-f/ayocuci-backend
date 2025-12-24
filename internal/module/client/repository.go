package client

import "gorm.io/gorm"

type Repository interface {
	Create(client *Client) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(client *Client) error {
	return r.db.Create(client).Error
}
