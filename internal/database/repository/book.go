package repository

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Book struct {
	Code      string `gorm:"unique"`
	Name      string
	Price     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDelete  soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
