package repository

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Book struct {
	Tag       string
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDelete  soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
