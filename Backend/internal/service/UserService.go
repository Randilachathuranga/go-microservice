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
	"strconv"
	"time"
)

type USerService struct {
	Repo repository.UserRepository
	Crep repository.CatalogRepository

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

func (s USerService) CrateProfile(id uint, input dto.ProfileInput) error {

	//update user
	_, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		return err
	}
	//create the address
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostalCode,
		UserId:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}
func (s USerService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s USerService) UpdateProfile(id uint, input dto.ProfileInput) error {
	// Find the existing user
	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}
	// Update the user fields
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	// Save the updated user
	_, err = s.Repo.UpdateUser(id, user)

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostalCode,
		UserId:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s USerService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	//find user
	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return "", errors.New("user not found")
	}
	if user.USerType == domain.SELLER {
		return "", errors.New("you are already a seller")
	}
	//update the user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phonenumber,
		USerType:  domain.SELLER,
	})
	if err != nil {
		return "", err
	}
	//generating token
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.USerType)
	if err != nil {
		return "", err
	}
	//create bank account information
	_, err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccountNumber: input.BankAccountNumber,
		SwiftCode:         input.SwiftCode,
		PaymentType:       input.PaymentType,
		UserId:            id,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
func (s USerService) FindCart(id uint) ([]domain.Cart, error) {
	cartitems, err := s.Repo.FindCartItems(id)
	return cartitems, err
}
func (s USerService) CreateCart(input dto.CreateCartRequest, u domain.User) (domain.Cart, error) {
	// Check if the cart item already exists
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)

	// Validate product ID
	if input.ProductId == 0 {
		return domain.Cart{}, errors.New("invalid product id")
	}

	if cart.ID > 0 {
		// Update or delete existing cart item
		if input.Quantity < 1 {
			// Delete the cart item
			err := s.Repo.DeleteCartItems(cart.ID)
			if err != nil {
				return domain.Cart{}, errors.New("error deleting cart item")
			}
			return domain.Cart{}, nil // Item deleted, return empty cart
		} else {
			// Update the quantity
			cart.Quantity = input.Quantity
			err := s.Repo.CreateCart(cart) // assuming CreateCart upserts
			if err != nil {
				return domain.Cart{}, errors.New("error updating cart")
			}
			return cart, nil
		}
	} else {
		//check if product exist
		product, err := s.Crep.FindProductById(int(input.ProductId))
		if product.ID < 1 {
			return domain.Cart{}, errors.New("Not found product")
		}
		// Create new cart item
		newCart := domain.Cart{
			UserID:    u.ID,
			ProductID: input.ProductId,
			Name:      product.Name,
			ImageURL:  product.ImageUrl,
			Quantity:  input.Quantity,
			Price:     float32(product.Price),
			SellerID:  product.UserID,
		}
		err = s.Repo.CreateCart(newCart)
		if err != nil {
			return domain.Cart{}, errors.New("error creating new cart item")
		}
		return newCart, nil
	}
}

func (s USerService) CreateOrder(u domain.User) (int, error) {
	//find cart items for user
	cartItems, err := s.Repo.FindCartItems(u.ID)
	if err != nil {
		return 0, errors.New("error getting cart")
	}
	if len(cartItems) == 0 {
		return 0, errors.New("cart is empty and cannot create order")
	}
	//find success payment status
	paymentId := "PAY12345"
	txnID := "TNX12345"
	OrderRef, _ := helper.Randomnumbers(8)
	//create order with generated order number
	var Amount float32
	var OrderItem []domain.OrderItem

	for _, item := range cartItems {
		Amount += item.Price * float32(item.Quantity)
		OrderItem = append(OrderItem, domain.OrderItem{
			ProductID: strconv.Itoa(int(item.ProductID)),
			Qty:       strconv.Itoa(int(item.Quantity)),
			Price:     float64(item.Price),
			Name:      item.Name,
			ImageURL:  item.ImageURL,
			SellerId:  item.SellerID,
		})
	}

	order := domain.Order{
		UserId:         u.ID,
		PaymentId:      paymentId,
		TransactionId:  txnID,
		OrderRefNumber: uint(OrderRef),
		Amount:         float64(Amount),
		Items:          OrderItem,
	}

	err = s.Repo.CreateOrder(order)
	if err != nil {
		return 0, err
	}
	//send email to user

	// remove cart items from cart
	err = s.Repo.DeleteCartItems(u.ID)
	if err != nil {
		return 0, errors.New("error deleting cart")
	}
	//return order number
	return OrderRef, nil
}
func (s USerService) GetOrders(u domain.User) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, errors.New("error getting orders")
	}

	return orders, nil
}

func (s USerService) GetorderbyID(id uint, uId uint) (domain.Order, error) {
	order, err := s.Repo.FindOrderByid(id, uId)
	if err != nil {
		return order, errors.New("error getting order")
	}
	return order, nil
}
