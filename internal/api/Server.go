package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/Handlers"
)

func StartServer(config Config.AppConfig) {
	app := fiber.New() // declare a new Fiber app

	rh := &rest.RestHandler{
		App: app,
	}

	SetupRoutes(rh)

	fmt.Printf("Server starting on port %s\n", config.ServerPort)
	if err := app.Listen(config.ServerPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func SetupRoutes(rh *rest.RestHandler) {
	//user handler
	Handlers.SetUpuserRoutes(rh)
	//transaction
	//catalog
}
