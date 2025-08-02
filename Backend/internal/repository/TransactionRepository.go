package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

// Interface
type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

// Implementation
type transactionStorage struct {
	db *gorm.DB
}

func (r *transactionStorage) CreatePayment(payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}
	return r.db.Create(payment).Error
}

func (r *transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	if uId == 0 {
		return nil, errors.New("invalid user ID")
	}

	var orders []domain.OrderItem
	err := r.db.Where("seller_id = ?", uId).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	if uId == 0 {
		return dto.SellerOrderDetails{}, errors.New("invalid user ID")
	}
	if id == 0 {
		return dto.SellerOrderDetails{}, errors.New("invalid order ID")
	}

	var details dto.SellerOrderDetails
	// More comprehensive query with proper joins
	query := `
		SELECT 
			o.order_ref_number,
			o.status as order_status,
			o.created_at,
			oi.id as order_item_id,
			oi.product_id,
			oi.name,
			oi.image_url,
			oi.price,
			oi.qty as quantity,
			u.first_name as customer_name,
			u.email as customer_email,
			u.phone as customer_phone,
			COALESCE(a.address_line, '') as customer_address
		FROM order_items oi
		INNER JOIN orders o ON oi.order_id = o.id
		INNER JOIN users u ON o.user_id = u.id
		LEFT JOIN addresses a ON u.id = a.user_id AND a.is_default = true
		WHERE oi.seller_id = ? AND oi.id = ?
	`

	err := r.db.Raw(query, uId, id).Scan(&details).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SellerOrderDetails{}, errors.New("order not found")
		}
		return dto.SellerOrderDetails{}, err
	}

	return details, nil
}

// Constructor
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{
		db: db,
	}
}
