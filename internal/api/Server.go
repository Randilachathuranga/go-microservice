package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-ecommerce-app/Config"
)

func StartServer(config Config.AppConfig) {
	app := fiber.New() // declare a new Fiber app

	app.Get("/new", Handlecheck)
	app.Get("/", Handlecheck)

	fmt.Printf("Server starting on port %s\n", config.ServerPort)
	if err := app.Listen(config.ServerPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func Handlecheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "I am Okk",
	})
}
