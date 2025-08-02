package domain

import "time"

type OrderItem struct {
	ID        int       `gorm:"primary_key" json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       string    `json:"qty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
