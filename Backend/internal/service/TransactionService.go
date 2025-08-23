package service

import (
	"errors"

	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"

	"github.com/stripe/stripe-go/v78"
)

// TransactionService defines the business logic for transactions.
type TransactionService interface {
	CreatePayment(payment *domain.Payment) error
	GetOrders(user domain.User) ([]domain.OrderItem, error)
	GetOrderDetails(user domain.User, orderID uint) (dto.SellerOrderDetails, error)
	GetActivePayment(userID uint) (*domain.Payment, error)
	StoreCreatedPayment(userID uint, ps *stripe.PaymentIntent, amount float32) error
	UpdatePayment(userID uint, status string, logs string) error
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
		payment.Status = domain.PaymentStatusInitial
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

// GetActivePayment retrieves the most recent active payment for a given user.
func (s *transactionService) GetActivePayment(userID uint) (*domain.Payment, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.Repo.FindPayment(userID)
}

// StoreCreatedPayment saves a Stripe PaymentIntent as a payment record.
func (s *transactionService) StoreCreatedPayment(userID uint, ps *stripe.PaymentIntent, amount float32) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}
	if ps == nil {
		return errors.New("payment intent cannot be nil")
	}

	payment := domain.Payment{
		UserId:    userID,
		Amount:    float64(amount),
		Status:    domain.PaymentStatusInitial,
		PaymentId: ps.ID,
		Response:  "",
	}

	return s.Repo.CreatePayment(&payment)
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

// UpdatePayment updates the payment status and response logs for a user's active payment
func (s *transactionService) UpdatePayment(userID uint, status string, logs string) error {
	if userID == 0 {
		return errors.New("invalid user id")
	}
	return s.Repo.UpdatePaymentStatus(userID, status, logs)
}
