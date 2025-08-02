package domain

import "time"

// this is use for database
const (
	SELLER = "seller"
	BUYER  = "buyer"
)

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	Code      uint      `json:"code"`
	Expiry    time.Time `json:"expiry"`
	Address   Address   `json:"address"` // relation is one to one
	Verified  bool      `json:"verified" gorm:"default:false"`
	USerType  string    `json:"user_type" gorm:"default:'buyer'"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
