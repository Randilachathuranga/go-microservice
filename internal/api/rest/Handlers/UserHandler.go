package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
)

type UserHandelr struct {
	//user service
	svc service.USerService
}

func SetUpuserRoutes(rh *rest.RestHandler) {

	app := rh.App

	//create an instance of user service and handler
	svc := service.USerService{
		Repo: repository.NewUserRepository(rh.DB),
		Auth: rh.Auth,
	}
	handler := UserHandelr{
		svc: svc,
	}

	pubRoutes := app.Group("/users")

	//Public endpoint
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)
	//private endpoint
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)
	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)

	pvtRoutes.Post("/cart", handler.AddtoCart)
	pvtRoutes.Get("/cart", handler.GetCart)
	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrder)

	pvtRoutes.Post("/become-seller", handler.BecomeaSeller)

}

// register a user
func (h *UserHandelr) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignUp{}
	err := ctx.BodyParser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	token, err := h.svc.Signup(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error on signup or missin token",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success register",
		"token":   token,
	})
}

func (h *UserHandelr) Login(ctx *fiber.Ctx) error {

	logininput := dto.UserLoging{}
	err := ctx.BodyParser(&logininput)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}

	token, err := h.svc.Login(logininput.Email, logininput.Password)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login",
		"token":   token,
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

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	log.Println(user)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile",
		"user":    user,
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
