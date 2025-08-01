package rest

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorMessage(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"message": message,
	})
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
}

func BadRequestError(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message,
	})
}

func SuccessResponse(ctx *fiber.Ctx, Msg string, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"Message": Msg,
		"data":    data,
	})
}
