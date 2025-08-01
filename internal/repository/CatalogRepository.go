package repository

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"gorm.io/gorm"
)

// Public interface
type CatalogRepository interface {
	CreateCategory(c *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id int) (*domain.Category, error)
	EditCategory(c *domain.Category) error
	DeleteCategory(c *domain.Category) error

	CreateProduct(c *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductById(id int) (*domain.Product, error)
	FindSellerProducts(id int) ([]*domain.Product, error)
	EditProduct(c *domain.Product) error
	DeleteProduct(c *domain.Product) error
}

// Private struct
type catalogRepository struct {
	db *gorm.DB
}

func (r *catalogRepository) CreateProduct(c *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) FindProducts() ([]*domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) FindProductById(id int) (*domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) FindSellerProducts(id int) ([]*domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) EditProduct(c *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) DeleteProduct(c *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *catalogRepository) CreateCategory(c *domain.Category) error {
	err := r.db.Create(c).Error
	fmt.Println(err)
	if err != nil {
		return errors.New("created successfully")
	}
	return nil
}

func (r *catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, errors.New("find all categories failed")
	}
	return categories, err
}

func (r *catalogRepository) FindCategoryById(id int) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *catalogRepository) EditCategory(c *domain.Category) error {
	err := r.db.Save(c).Error
	if err != nil {
		return errors.New("update category failed")
	}
	return nil
}

func (r *catalogRepository) DeleteCategory(c *domain.Category) error {
	err := r.db.Delete(c).Error
	if err != nil {
		return errors.New("delete category failed")
	}
	return nil
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
