package models

type EmailSignupRequest struct{
	Name string `validate:"required" json:"name"`
	Email string `validate:"required,email" json:"email"`
	PhoneNumber string `validate:"required,number,len=10,numeric" json:"phone_number"`
	Password string `validate:"required" json:"password"`
	ConfirmPassword string `validate:"required" json:"confirmpassword"`
}

type EmailLoginRequest struct {
	Email    string `form:"email" validate:"required,email" json:"email"`
	Password string `form:"password" validate:"required" json:"password"`
}

type SellerRegisterRequest struct {
	UserName    string `validate:"required" json:"name"`
	Password    string `validate:"required" json:"password"`
	Description string `validate:"required" json:"description"`
}

type SellerLoginRequest struct {
	UserID   uint   `json:"userid"`
	UserName string `form:"" username:"required" json:"username"`
	Password string `form:"password" validate:"required" json:"password"`
}
//
type EmailVerificationRequest struct{
	Email string `json:"email" validate:"required,email"`
	OTP  string  `json:"otp" validate:"required"`
}