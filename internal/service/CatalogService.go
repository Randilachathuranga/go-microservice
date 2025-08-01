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
