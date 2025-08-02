package domain

import "time"

type Cart struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	SellerID  uint      `json:"seller_id"`
	Price     float32   `json:"price"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
