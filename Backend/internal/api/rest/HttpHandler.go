package rest

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/pkg/Payment"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config Config.AppConfig
	Pc     payment.PaymentClient
}
