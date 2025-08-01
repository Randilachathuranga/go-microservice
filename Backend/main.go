package main

import (
	"github.com/gofiber/fiber/v2/log"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/api"
)

func main() {
	cfg, err := Config.SetupEnv()

	if err != nil {
		log.Fatal("Config file is not loaded ", err)
	}
	api.StartServer(cfg)
}
