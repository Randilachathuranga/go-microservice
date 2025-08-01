package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
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
	app.Get("/product", handler.GetAllPublicProducts)
	app.Get("/product/:id", handler.GetPublicProductByID)
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
	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	log.Println("user:", user)
	return rest.SuccessResponse(ctx, "Category created", nil)
}

func (h CatalogHandler) UpdateCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Category updated", nil)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Category deleted", nil)
}

func (h CatalogHandler) GetAllCategories(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "All categories fetched", nil)
}

func (h CatalogHandler) GetCategoryByID(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Category fetched by ID", nil)
}

// ----- Product Handlers (Public) -----
func (h CatalogHandler) GetAllPublicProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Public products fetched", nil)
}

func (h CatalogHandler) GetPublicProductByID(ctx *fiber.Ctx) error {
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
