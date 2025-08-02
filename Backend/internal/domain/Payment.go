package domain

import "time"

type Payment struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	TransactionId uint      `json:"transaction_id"`
	CustomerId    uint      `json:"customer_id"`
	PaymentId     uint      `json:"payment_id"`
	Status        string    `json:"status"`
	Response      string    `json:"response"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
