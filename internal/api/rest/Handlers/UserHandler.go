package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
)

type UserHandelr struct {
	//user service

}

func SetUpuserRoutes(rh *rest.RestHandler) {

	app := rh.App

	//create an instance of user service and handler

	handler := UserHandelr{}

	//Public endpoint
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	//private endpoint
	app.Get("/verify", handler.GetVerificationCode)
	app.Post("/verify", handler.Verify)
	app.Post("/Profile", handler.CreateProfile)
	app.Get("/Profile", handler.GetProfile)

	app.Post("/cart", handler.AddtoCart)
	app.Get("/cart", handler.GetCart)
	app.Get("/order", handler.GetOrders)
	app.Get("/order/:id", handler.GetOrder)

	app.Post("/become-seller", handler.BecomeaSeller)

}

func (h *UserHandelr) Register(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "register",
	})
}
func (h *UserHandelr) Login(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login",
	})
}
func (h *UserHandelr) Verify(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Verify",
	})
}
func (h *UserHandelr) GetVerificationCode(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get verification code",
	})
}
func (h *UserHandelr) CreateProfile(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Create profile",
	})
}
func (h *UserHandelr) GetProfile(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile",
	})
}
func (h *UserHandelr) AddtoCart(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Add to cart",
	})
}
func (h *UserHandelr) GetCart(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get cart",
	})
}
func (h *UserHandelr) CreateCart(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Create cart",
	})
}
func (h *UserHandelr) GetOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id") // get the :id parameter from URL
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get a order",
		"orderID": orderID,
	})
}
func (h *UserHandelr) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get orderss",
	})
}
func (h *UserHandelr) BecomeaSeller(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "become a seller",
	})
}
