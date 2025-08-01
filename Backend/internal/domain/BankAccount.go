package domain

import "time"

type BankAccount struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	UserId            uint      `json:"user-id"`
	BankAccountNumber string    `json:"bank_account_number"`
	SwiftCode         string    `json:"swift_code"`
	PaymentType       string    `json:"payment_type"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
