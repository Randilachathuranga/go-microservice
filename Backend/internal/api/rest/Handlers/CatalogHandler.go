package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/domain"
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
	app.Get("/product", handler.GetProducts)
	app.Get("/product/:id", handler.GetProduct)
	app.Get("/categories", handler.GetAllCategories)
	app.Get("/categories/:id", handler.GetCategoryByID)

	// Seller (private) endpoints
	selroutes := app.Group("/seller", rh.Auth.AuthorizeSeller)

	// Categories
	selroutes.Post("/categories", handler.CreateCategories)
	selroutes.Patch("/categories/:id", handler.EditCategory)
	selroutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	selroutes.Post("/products", handler.CreateProducts)
	selroutes.Get("/products", handler.GetProducts)
	selroutes.Get("/products/:id", handler.GetProduct)
	selroutes.Put("/products/:id", handler.EditProduct)
	selroutes.Patch("/products/:id", handler.UpdateStock)
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

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
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

// ----- Product Handlers -----

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	// Check if this is a seller endpoint by looking at the path
	if ctx.Path() == "/seller/products" {
		// This is seller products endpoint - get products for current seller
		_, err := h.svc.Auth.GetCurrentUser(ctx)
		if err != nil {
			return rest.BadRequestError(ctx, "Unauthorized")
		}

		products, err := h.svc.GetSellerProducts()
		if err != nil {
			return rest.InternalError(ctx, err)
		}
		return rest.SuccessResponse(ctx, "Seller products fetched", products)
	}

	// This is public products endpoint
	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Products fetched", products)
}

func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid product ID")
	}

	product, err := h.svc.GetProductById(id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Product fetched by ID", product)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	req := dto.CreateProductRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "Cannot create product")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.BadRequestError(ctx, "Unauthorized")
	}

	if err := h.svc.CreateProduct(req, user); err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Product created successfully", req)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid product ID")
	}

	var req dto.CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "Invalid request body")
	}

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.BadRequestError(ctx, "Unauthorized")
	}

	updatedProduct, err := h.svc.EditProduct(id, req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Product updated successfully", updatedProduct)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid product ID")
	}
	// Parse stock from request body
	var req struct {
		Stock int `json:"stock"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestError(ctx, "Invalid request body")
	}
	// Get current user (auth required)
	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.BadRequestError(ctx, "Unauthorized")
	}
	// Create product with updated stock
	product := domain.Product{
		ID:     uint(id),
		Stock:  uint(req.Stock),
		UserID: user.ID,
	}
	// Call service to update stock
	updatedProduct, err := h.svc.UpdateProductStock(product)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Product stock updated successfully", updatedProduct)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	// Parse product ID from URL
	_, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return rest.BadRequestError(ctx, "Invalid product ID")
	}

	// Get current user
	_, err = h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.BadRequestError(ctx, "Unauthorized")
	}

	// Call service to delete the product
	if err := h.svc.DeleteProduct(ctx); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "Product deleted successfully", nil)
}
