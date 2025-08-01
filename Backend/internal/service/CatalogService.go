package service

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config Config.AppConfig
}

func (c *CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := c.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: int(input.DisplayOrder),
		ParentID:     input.ParentId,
	})
	return err
}

func (c *CatalogService) EditCategory(id int, request dto.CreateCategoryRequest) (*domain.Category, error) {
	category, err := c.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = request.Name
	category.ImageUrl = request.ImageUrl
	category.DisplayOrder = int(request.DisplayOrder)
	category.ParentID = request.ParentId

	err = c.Repo.EditCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *CatalogService) DeleteCategory(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid category ID")
	}

	category, err := c.Repo.FindCategoryById(id)
	if err != nil {
		return rest.BadRequestError(ctx, "Category not found")
	}

	err = c.Repo.DeleteCategory(category)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "Category deleted successfully", nil)
}

func (c *CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := c.Repo.FindCategories()
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (c *CatalogService) GetCategory(id int) (*domain.Category, error) {
	cat, err := c.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("Error")
	}
	return cat, nil
}

// products
func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		ImageUrl:    input.ImageUrl,
		UserID:      user.ID,
		Stock:       input.Stock,
	})
	return err
}

func (s *CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.UserID != user.ID {
		return nil, errors.New("user not authorized to edit")
	}

	// Update fields
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.CategoryId = input.CategoryId
	product.ImageUrl = input.ImageUrl
	// Save changes
	if err := s.Repo.EditProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *CatalogService) DeleteProduct(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid category ID")
	}
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return rest.BadRequestError(ctx, "Product not found")
	}
	err = s.Repo.DeleteProduct(product)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Product deleted successfully", nil)
}

func (s *CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, err
	}
	return products, err
}
func (s *CatalogService) GetProductById(id int) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("Product not found")
	}
	return product, nil
}

func (s *CatalogService) GetSellerProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, err
	}
	return products, err
}

func (s *CatalogService) UpdateProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(int(e.ID))
	if err != nil {
		return nil, errors.New("product not found")
	}
	product.Stock = e.Stock
	if err := s.Repo.EditProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}
