package Handlers

import (
	"strconv"

	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc  service.TransactionService
	auth helper.Auth
}

func InitializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	repo := repository.NewTransactionRepository(db)
	return service.NewTransactionService(repo, auth)
}

func SetTransactionRoutes(as *rest.RestHandler) {
	app := as.App
	svc := InitializeTransactionService(as.DB, as.Auth)

	handler := TransactionHandler{
		svc:  svc,
		auth: as.Auth,
	}

	secRoute := app.Group("/", as.Auth.Authorize)
	secRoute.Post("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.Authorize)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	// Get current user from context
	user, err := h.auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"error":   err.Error(),
		})
	}

	// Parse payment request
	var paymentRequest struct {
		Amount        float64 `json:"amount" validate:"required,min=0.01"`
		CaptureMethod string  `json:"capture_method" validate:"required"`
		CustomerId    uint    `json:"customer_id" validate:"required"`
		TransactionId uint    `json:"transaction_id" validate:"required"`
	}

	if err := ctx.BodyParser(&paymentRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if paymentRequest.Amount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "amount must be greater than 0",
		})
	}

	if paymentRequest.CaptureMethod == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "capture method is required",
		})
	}

	// Create payment object
	payment := &domain.Payment{
		UserId:        user.ID,
		Amount:        paymentRequest.Amount,
		CaptureMethod: paymentRequest.CaptureMethod,
		CustomerId:    paymentRequest.CustomerId,
		TransactionId: paymentRequest.TransactionId,
		Status:        "pending",
		Response:      "",
	}

	// Process payment through service
	if err := h.svc.CreatePayment(payment); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to process payment",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "payment processed successfully",
		"payment_id": payment.ID,
		"status":     payment.Status,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	// Get current user from context
	user, err := h.auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"error":   err.Error(),
		})
	}

	// Get orders from service
	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve orders",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "orders retrieved successfully",
		"data":    orders,
		"count":   len(orders),
	})
}
func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	// Get current user from context
	user, err := h.auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"error":   err.Error(),
		})
	}

	// Parse order ID from URL params
	orderIdParam := ctx.Params("id")
	if orderIdParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "order ID is required",
		})
	}

	orderID, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order ID format",
			"error":   err.Error(),
		})
	}

	// Get order details from service
	orderDetails, err := h.svc.GetOrderDetails(user, uint(orderID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve order details",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "order details retrieved successfully",
		"data":    orderDetails,
	})
}
