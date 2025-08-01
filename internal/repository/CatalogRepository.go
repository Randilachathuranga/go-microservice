package repository

import (
	"gorm.io/gorm"
)

// Public interface
type CatalogRepository interface {
}

// Private struct
type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
