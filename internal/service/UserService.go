package service

import (
	"fmt"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"log"
)

type USerService struct {
	Repo repository.UserRepository
}

// input any == input interface == any parameters
func (s USerService) Signup(input dto.UserSignUp) (string, error) {
	log.Println(input)
	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: input.Password,
		Phone:    input.Phone,
	})

	//generate token
	log.Println(user)
	userInfo := fmt.Sprintln("New user: ", user)

	return userInfo, err
}
func (s USerService) Login(email string, password string) (string, error) {

	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	log.Println(user)
	return user.Email, err
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
