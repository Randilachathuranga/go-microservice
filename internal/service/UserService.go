package service

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"log"
)

type USerService struct {
	Repo repository.UserRepository
	Auth helper.Auth
}

// input any == input interface == any parameters
func (s USerService) Signup(input dto.UserSignUp) (string, error) {
	hpassword, err := s.Auth.CreateHashPassword(input.Password)
	if err != nil {
		return "", err
	}
	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Phone:    input.Phone,
		Password: hpassword,
	})
	//generate token
	log.Println(user)
	return s.Auth.GenerateToken(user.ID, user.Email, user.USerType)
}

// login
func (s USerService) Login(email string, password string) (string, error) {
	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		return "", err // user not found
	}

	// ✅ THIS LINE MUST RUN
	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err // invalid password
	}

	// ✅ Only generate token if password is valid
	return s.Auth.GenerateToken(user.ID, user.Email, user.USerType)
}

func (s USerService) findbyEmail(email string) (domain.User, error) {
	//perform sone db operation
	//business logic
	user, err := s.Repo.FindUserByEmail(email)
	return user, err
}

func (s USerService) GetverificationCode(e domain.User) (int, error) {
	return 0, nil
}

func (s USerService) VerifyCode(id uint, code uint) error {
	return nil
}
func (s USerService) CeateProfile(id uint, input any) error {
	return nil
}
func (s USerService) GetProfile(id uint) (*domain.User, error) {
	return nil, nil
}
func (s USerService) UpdateProfile(id uint, input any) error {
	return nil
}

func (s USerService) BecomeSeller(id uint, input any) (string, error) {
	return "", nil
}
func (s USerService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}
func (s USerService) CreateCart(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}
func (s USerService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}
func (s USerService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s USerService) GetorderbyID(id uint, Uid uint) (interface{}, error) {
	return nil, nil
}
