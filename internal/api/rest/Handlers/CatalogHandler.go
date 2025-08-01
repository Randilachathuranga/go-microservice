package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"strconv"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := CatalogHandler{
		svc: svc,
	}

	// Public endpoints
	app.Get("/product", handler.GetAllProducts)
	app.Get("/product/:id", handler.GetProductByID)
	app.Get("/categories", handler.GetAllCategories)
	app.Get("/categories/:id", handler.GetCategoryByID)

	// Seller (private) endpoints
	selroutes := app.Group("/seller", rh.Auth.AuthorizeSeller)

	// Categories
	selroutes.Post("/categories", handler.CreateCategories)
	selroutes.Patch("/categories/:id", handler.UpdateCategory)
	selroutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	selroutes.Post("/products", handler.CreateProduct)
	selroutes.Get("/products", handler.GetAllSellerProducts)
	selroutes.Get("/products/:id", handler.GetSellerProductByID)
	selroutes.Put("/products/:id", handler.ReplaceProduct)
	selroutes.Patch("/products/:id", handler.UpdateProduct)
	selroutes.Delete("/products/:id", handler.DeleteProduct)
}

// ----- Category Handlers -----

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "Cannot create category")
	}

	err = h.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Category created successfully", nil)
}

func (h CatalogHandler) UpdateCategory(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid category ID")
	}

	var req dto.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "Invalid request body")
	}

	updatedCat, err := h.svc.EditCategory(id, req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "Category updated successfully", updatedCat)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return h.svc.DeleteCategory(ctx)
}

func (h CatalogHandler) GetAllCategories(ctx *fiber.Ctx) error {
	cats, err := h.svc.GetCategories()
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "All categories fetched successfully", cats)
}

func (h CatalogHandler) GetCategoryByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid category ID")
	}

	cat, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Category fetched by ID", cat)
}

// ----- Product Handlers (Public) -----

func (h CatalogHandler) GetAllProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Public products fetched", nil)
}

func (h CatalogHandler) GetProductByID(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Public product fetched by ID", nil)
}

// ----- Product Handlers (Seller) -----

func (h CatalogHandler) CreateProduct(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, "Product created", nil)
}

func (h CatalogHandler) GetAllSellerProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Seller products fetched", nil)
}

func (h CatalogHandler) GetSellerProductByID(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Seller product fetched by ID", nil)
}

func (h CatalogHandler) ReplaceProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Product replaced", nil)
}

func (h CatalogHandler) UpdateProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Product updated", nil)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Product deleted", nil)
}
