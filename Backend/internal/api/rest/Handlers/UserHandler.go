package Handlers

import (
	"fmt"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandelr struct {
	//user service
	svc service.USerService
}

func SetUpuserRoutes(rh *rest.RestHandler) {

	app := rh.App

	//create an instance of user service and handler
	svc := service.USerService{
		Repo:   repository.NewUserRepository(rh.DB),
		Crep:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
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
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Post("/cart", handler.AddtoCart)
	pvtRoutes.Get("/cart", handler.GetCart)

	pvtRoutes.Post("/order", handler.CreateOrder)
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

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	//request
	var req dto.VerificationCodeInput
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)

	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Verified Successfully",
	})
}

func (h *UserHandelr) GetVerificationCode(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	//create verfification code and update to user profile
	err := h.svc.GetverificationCode(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error on get verification code",
		})

	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get verification code",
	})
}

func (h *UserHandelr) CreateProfile(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	err = h.svc.CrateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success create profile",
	})
}
func (h *UserHandelr) GetProfile(ctx *fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	log.Println(user)

	//
	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile",
		"Profile": profile,
	})
}

func (h *UserHandelr) UpdateProfile(ctx *fiber.Ctx) error {
	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	log.Println(user)

	req := dto.ProfileInput{}
	err = ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	err = h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

func (h *UserHandelr) AddtoCart(ctx *fiber.Ctx) error {

	req := dto.CreateCartRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs products and quantity",
		})
	}
	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	fmt.Println(user)

	//call user service and perform create cart
	cartitems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessResponse(ctx, "Cart Items created successfully", cartitems)
}
func (h *UserHandelr) GetCart(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	cart, err := h.svc.FindCart(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return rest.SuccessResponse(ctx, "Cart Details", cart)
}
func (h *UserHandelr) CreateOrder(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	orderRef, err := h.svc.CreateOrder(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Order Created",
		"OrderRef": orderRef,
	})
}
func (h *UserHandelr) GetOrder(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	orderid, _ := strconv.Atoi(ctx.Params("id"))

	order, err := h.svc.GetorderbyID(uint(orderid), user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get a order",
		"orderID": order,
	})
}
func (h *UserHandelr) GetOrders(ctx *fiber.Ctx) error {

	user, _ := h.svc.Auth.GetCurrentUser(ctx)
	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get orderss",
		"orders":  orders,
	})
}
func (h *UserHandelr) BecomeaSeller(ctx *fiber.Ctx) error {

	user, err := h.svc.Auth.GetCurrentUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	//
	req := dto.SellerInput{}
	err = ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	token, err := h.svc.BecomeSeller(user.ID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "fail to become seller",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Become seller",
		"token":   token,
	})
}
