package Handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"

	"github.com/gofiber/fiber/v2"
)

// SetupCompatibilityRoutes adds frontend-friendly endpoints that map to internal handlers
func SetupCompatibilityRoutes(rh *rest.RestHandler) {
	app := rh.App

	// health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	// collect-payment -> uses buyer MakePayment (POST expected by frontend)
	app.Post("/collect-payment", func(c *fiber.Ctx) error {
		// reuse existing transaction handler
		svc := initializeTransactionService(rh.DB, rh.Auth)
		h := TransactionHandler{Svc: svc, PaymentClient: rh.Pc, UserSvc: service.USerService{Repo: repository.NewUserRepository(rh.DB), Auth: rh.Auth, Crep: repository.NewCatalogRepository(rh.DB), Config: rh.Config}, Config: rh.Config}
		return h.MakePayment(c)
	})

	// order -> create order via users handler
	app.Post("/order", func(c *fiber.Ctx) error {
		// reuse user handler
		svc := service.USerService{Repo: repository.NewUserRepository(rh.DB), Crep: repository.NewCatalogRepository(rh.DB), Auth: rh.Auth, Config: rh.Config}
		h := UserHandelr{svc: svc}
		return h.CreateOrder(c)
	})
}
