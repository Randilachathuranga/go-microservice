package service

import (
	"errors"
	"fmt"
	"go-ecommerce-app/Config"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/Notification"
	"log"
	"time"
)

type USerService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
	Config Config.AppConfig
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

func (s USerService) isverifiedUser(id uint) bool {
	currentUser, err := s.Repo.FindUserById(id)
	return err == nil && currentUser.Verified
}

func (s USerService) GetverificationCode(e domain.User) error {
	// if user already verified
	if s.isverifiedUser(e.ID) {
		return errors.New("verification code already used")
	}

	//generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return err
	}
	//update user
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   uint(code),
	}
	_, err = s.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return err
	}

	user, errs := s.Repo.FindUserById(e.ID)
	if err != nil {
		return errs
	}
	//send SMS
	notificationClient := Notification.NewNotificationClient(s.Config)

	msg := fmt.Sprintf("Your verification code is: %s", code)

	err = notificationClient.SendSMS(user.Phone, msg)
	if err != nil {
		return err
	}

	//return the verification code
	return nil
}

func (s USerService) VerifyCode(id uint, code uint) error {
	if s.isverifiedUser(id) {
		log.Println("verified ...")
		return errors.New("verification code already used")
	}

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Code != code {
		return errors.New("verification code does not match")
	}
	if time.Now().After(user.Expiry) {
		return errors.New("verification code expired")
	}
	updater := domain.User{
		Verified: true,
	}
	_, err = s.Repo.UpdateUser(id, updater)
	if err != nil {
		return err
	}

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
