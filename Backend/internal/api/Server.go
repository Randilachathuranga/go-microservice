package api

import (
	"fmt"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/Handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config Config.AppConfig) {
	app := fiber.New() // declare a new Fiber app

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // React dev server port
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database", err)
	}
	fmt.Println("database connection established", db)

	//run the migration
	err = db.AutoMigrate(
		&domain.User{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Address{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)
	if err != nil {
		log.Warnf("Migration warning: %s", err.Error())
		fmt.Println("Migration completed with warnings - continuing...")
	} else {
		fmt.Println("Database migrated successfully")
	}

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
	Handlers.SetTransactionRoutes(rh)
	//catalog
	Handlers.SetCatalogRoutes(rh)
}
