package Handlers

import (
	"encoding/json"
	"errors"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	payment "go-ecommerce-app/pkg/Payment"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	Svc           service.TransactionService
	UserSvc       service.USerService
	PaymentClient payment.PaymentClient
	Config        Config.AppConfig
}

// initialize service
func initializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.NewTransactionService(repository.NewTransactionRepository(db), auth)
}

// setup routes
func SetupTransactionRoutes(as *rest.RestHandler) {
	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)

	useSvc := service.USerService{
		Repo:   repository.NewUserRepository(as.DB),
		Crep:   repository.NewCatalogRepository(as.DB),
		Auth:   as.Auth,
		Config: as.Config,
	}

	handler := TransactionHandler{
		Svc:           svc,
		PaymentClient: as.Pc,
		UserSvc:       useSvc,
		Config:        as.Config,
	}

	// buyer routes
	secRoute := app.Group("/buyer", as.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)
	secRoute.Get("/verify", handler.VerifyPayment)

	// seller routes
	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

// ========== HANDLERS ==========

// create payment
func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	// 1. grab authorized user
	user, err := h.UserSvc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	// public key isn't present in AppConfig; use Stripe public key if needed from Config (field StripeSecret exists)
	pubKey := h.Config.StripeSecret

	// 2. check active payment session
	activePayment, err := h.Svc.GetActivePayment(user.ID)
	if activePayment != nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "create payment",
			"pubKey":  pubKey,
			"secret":  activePayment.PaymentId,
		})
	}

	// 3. get cart total
	_, amount, err := h.UserSvc.FindCart(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// generate orderId
	orderId, err := helper.Randomnumbers(8)
	if err != nil {
		return rest.InternalError(ctx, errors.New("error generating order id"))
	}

	// 4. create new payment session
	paymentResult, err := h.PaymentClient.CreatePayment(float64(amount), user.ID, orderId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// 5. store payment in db
	// Store session using service method signature
	err = h.Svc.StoreCreatedPayment(user.ID, paymentResult, amount)
	if err != nil {
		return ctx.Status(400).JSON(err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "create payment",
		"pubKey":  pubKey,
		"secret":  paymentResult.ID,
	})
}

// verify payment
func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user, err := h.UserSvc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// check active payment
	activePayment, err := h.Svc.GetActivePayment(user.ID)
	if err != nil || activePayment.ID == 0 {
		return ctx.Status(400).JSON(errors.New("no active payment exist"))
	}

	// fetch status from provider
	paymentRes, err := h.PaymentClient.GetPaymentStatus(activePayment.PaymentId)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"

	// success → create order
	if paymentRes.Status == "succeeded" {
		paymentStatus = "success"
		// Create order using the user service; CreateOrder expects a domain.User
		_, err = h.UserSvc.CreateOrder(user)
	}

	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// update status
	h.Svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":  "verify payment",
		"response": paymentRes,
	})
}

// get all orders (seller)
func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

// get order details (seller)
func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}
