package domain

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryId  uint      `json:"categoryId"`
	ImageUrl    string    `json:"image_url"`
	Price       uint      `json:"price"`
	UserID      uint      `json:"user_id"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
