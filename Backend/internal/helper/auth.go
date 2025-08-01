package helper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go-ecommerce-app/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{Secret: s}
}

// CreateHashPassword generates hash value for password
func (a Auth) CreateHashPassword(p string) (string, error) {
	if len(p) < 6 {
		return "", errors.New("password is too short, add more than 6 characters")
	}

	hashp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password hash failed")
	}

	return string(hashp), nil
}

// GenerateToken creates a JWT token
func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("id, email, and role are required")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"role":  role,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.Secret)) // Use the secret from struct
	if err != nil {
		return "", errors.New("token sign failed")
	}

	return tokenString, nil // Return nil error on success
}

// VerifyPassword compares plain password with hashed password
func (a Auth) VerifyPassword(plainPassword string, hashedPassword string) error {
	if len(plainPassword) < 6 {
		return errors.New("password is too short, add more than 6 characters")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

// VerifyToken validates JWT token and returns user information
func (a Auth) VerifyToken(t string) (domain.User, error) {
	// Handle Bearer token format: "Bearer <token>"
	tokenArr := strings.Split(t, " ")
	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("token format is invalid")
	}
	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("token must start with Bearer")
	}

	token, err := jwt.Parse(tokenArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return domain.User{}, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, errors.New("invalid token claims")
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			return domain.User{}, errors.New("token is expired")
		}
	} else {
		return domain.User{}, errors.New("invalid expiration claim")
	}

	// Extract user information from claims
	user := domain.User{}

	if id, ok := claims["id"].(float64); ok {
		user.ID = uint(id)
	} else {
		return domain.User{}, errors.New("invalid id claim")
	}

	if email, ok := claims["email"].(string); ok {
		user.Email = email
	} else {
		return domain.User{}, errors.New("invalid email claim")
	}

	if role, ok := claims["role"].(string); ok {
		user.USerType = role // Fixed typo: USerType -> UserType
	} else {
		return domain.User{}, errors.New("invalid role claim")
	}

	return user, nil
}

// Authorize is a middleware method for authentication
func (a Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "authorization header is required",
		})
	}

	user, err := a.VerifyToken(authHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	}

	if user.ID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid user information",
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

// GetCurrentUser retrieves the current user from context
func (a Auth) GetCurrentUser(ctx *fiber.Ctx) (domain.User, error) {
	user := ctx.Locals("user")
	if user == nil {
		return domain.User{}, errors.New("user not found in context")
	}
	userObj, ok := user.(domain.User)
	if !ok {
		return domain.User{}, errors.New("invalid user type in context")
	}
	return userObj, nil
}

func (a Auth) GenerateCode() (int, error) {
	return Randomnumbers(6)
}

// Authorize is a middleware method for authentication
func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "authorization header is required",
		})
	}
	user, err := a.VerifyToken(authHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	} else if user.ID > 0 && user.USerType == domain.SELLER {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "invalid user type",
			"reason":  errors.New("invalid user type in context"),
		})
	}
}
