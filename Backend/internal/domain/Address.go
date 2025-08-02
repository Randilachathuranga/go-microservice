package domain

import "time"

type Address struct {
	ID           uint      `gorm:"primary_key;auto_increment"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2"`
	City         string    `json:"city"`
	PostCode     string    `json:"post_code"`
	Country      string    `json:"country"`
	UserId       uint      `json:"user_id"` //one to one
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
