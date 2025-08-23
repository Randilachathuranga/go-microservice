package Config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort   string
	Dsn          string
	AppSecret    string
	AccountSid   string
	AuthToken    string
	Fromphone    string
	StripeSecret string
	PubKey       string
	SuccessUrl   string
	CancelUrl    string
}

func SetupEnv() (AppConfig, error) {

	// Load .env file by default
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// Server Port
	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("environment variable HTTP_PORT is not set")
	}
	fmt.Println("HTTP_PORT is:", httpPort)
	if httpPort[0] != ':' {
		httpPort = ":" + httpPort
	}

	// Database DSN
	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("environment variable DSN is not set")
	}
	fmt.Println("Dsn is:", Dsn)

	// App Secret
	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("environment variable APP_SECRET is not set")
	}

	// Twilio / Messaging
	accountSid := os.Getenv("ACCOUNT_SID")
	if len(accountSid) < 1 {
		accountSid = "default_account_sid"
	}
	fmt.Println("AccountSid is:", accountSid)

	authToken := os.Getenv("AUTH_TOKEN")
	if len(authToken) < 1 {
		authToken = "default_auth_token"
	}
	fmt.Println("AuthToken is:", authToken)

	fromphone := os.Getenv("FROM_PHONE")
	if len(fromphone) < 1 {
		fromphone = "default_phone"
	}
	fmt.Println("Fromphone is:", fromphone)

	// Stripe Configuration
	stripeSecret := os.Getenv("STRIPE_SECRET")
	if len(stripeSecret) < 1 {
		stripeSecret = "default_stripe_secret"
	}

	pubKey := os.Getenv("PUB_KEY")
	// leave pubKey empty if not provided; frontend may need it

	successUrl := os.Getenv("SUCCESS_URL")
	if len(successUrl) < 1 {
		successUrl = "http://localhost:3000/success"
	}

	cancelUrl := os.Getenv("CANCEL_URL")
	if len(cancelUrl) < 1 {
		cancelUrl = "http://localhost:3000/cancel"
	}

	return AppConfig{
		ServerPort:   httpPort,
		Dsn:          Dsn,
		AppSecret:    appSecret,
		AccountSid:   accountSid,
		AuthToken:    authToken,
		Fromphone:    fromphone,
		StripeSecret: stripeSecret,
		PubKey:       pubKey,
		SuccessUrl:   successUrl,
		CancelUrl:    cancelUrl,
	}, nil
}
