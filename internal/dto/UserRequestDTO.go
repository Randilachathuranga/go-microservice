package dto

type UserLoging struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUp struct {
	UserLoging
	Phone string `json:"phone"`
}

type VerificationCodeInput struct {
	Code uint `json:"code"`
}

type SellerInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Phonenumber       string `json:"phonenumber"`
	BankAccountNumber string `json:"bank_account_number"`
	SwiftCode         string `json:"swift_code"`
	PaymentType       string `json:"payment_type"`
}
