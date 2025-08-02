package service

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

// TransactionService defines the business logic for transactions.
type TransactionService interface {
	CreatePayment(payment *domain.Payment) error
	GetOrders(user domain.User) ([]domain.OrderItem, error)
	GetOrderDetails(user domain.User, orderID uint) (dto.SellerOrderDetails, error)
}

// transactionService is the implementation of the TransactionService interface.
type transactionService struct {
	Repo repository.TransactionRepository
	Auth helper.Auth
}

// NewTransactionService returns a new instance of TransactionService.
func NewTransactionService(repo repository.TransactionRepository, auth helper.Auth) TransactionService {
	return &transactionService{
		Repo: repo,
		Auth: auth,
	}
}

// CreatePayment handles saving a payment to the database.
func (s *transactionService) CreatePayment(payment *domain.Payment) error {
	// Validate payment data
	if payment.Amount <= 0 {
		return errors.New("payment amount must be greater than 0")
	}
	if payment.CaptureMethod == "" {
		return errors.New("capture method is required")
	}
	if payment.UserId == 0 {
		return errors.New("user ID is required")
	}

	// Set default status if not provided
	if payment.Status == "" {
		payment.Status = "pending"
	}

	return s.Repo.CreatePayment(payment)
}

// GetOrders retrieves all orders for the given user.
func (s *transactionService) GetOrders(user domain.User) ([]domain.OrderItem, error) {
	if user.ID == 0 {
		return nil, errors.New("invalid user ID")
	}

	orders, err := s.Repo.FindOrders(user.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderDetails retrieves a specific order by ID for the given user.
func (s *transactionService) GetOrderDetails(user domain.User, orderID uint) (dto.SellerOrderDetails, error) {
	if user.ID == 0 {
		return dto.SellerOrderDetails{}, errors.New("invalid user ID")
	}
	if orderID == 0 {
		return dto.SellerOrderDetails{}, errors.New("invalid order ID")
	}

	order, err := s.Repo.FindOrderById(user.ID, orderID)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}

	return order, nil
}
