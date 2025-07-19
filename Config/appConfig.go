package Config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ServerPort string
	Dsn        string
	AppSecret  string
}

func SetupEnv() (cfg AppConfig, err error) {

	// Load .env file by default (you can set APP_ENV=prod to skip this)
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("Environment variable HTTP_PORT is not set")
	}
	fmt.Println("HTTP_PORT is:", httpPort)

	// Fix: Ensure port format is correct for Fiber
	if httpPort[0] != ':' {
		httpPort = ":" + httpPort
	}

	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("Environment variable DSN is not set")
	}
	fmt.Println("Dsn is:", Dsn)

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("Environment variable APP_SECRET is not set")
	}

	return AppConfig{httpPort, Dsn, appSecret}, nil
}
