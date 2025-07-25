package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/Handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config Config.AppConfig) {
	app := fiber.New() // declare a new Fiber app

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database", err)
	}
	fmt.Println("database connection established", db)

	//run the migration
	db.AutoMigrate(&domain.User{})

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
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
