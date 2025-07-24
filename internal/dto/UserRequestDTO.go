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
