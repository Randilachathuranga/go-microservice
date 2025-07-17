package service

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"log"
)

type USerService struct {
}

// input any == input interface == any parameters
func (s USerService) Signup(input dto.UserSignUp) (string, error) {
	log.Println(input)
	return "this is my token", nil
}

func (s USerService) findbyEmail(email string) (domain.User, error) {
	//perform sone db operation
	//business logic
	return domain.User{}, nil
}

func (s USerService) login(input any) (string, error) {
	return "", nil
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
