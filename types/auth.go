package types

type SendOtpEmailBody struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailBody struct {
	Otp   string `json:"otp" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type RegisterBody struct {
	FirstName   string `json:"first_name" binding:"required,min=3"`
	LastName    string `json:"last_name" binding:"required,min=3"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginBody struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required,min=8"`
}

type RefreshTokenBody struct {
	Token string `json:"refresh_token" binding:"required"`
}
