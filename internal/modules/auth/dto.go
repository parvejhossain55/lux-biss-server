package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/parvej/luxbiss_server/pkg/validators"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128,secure_password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	OTP             string `json:"otp" validate:"required,len=6"`
	Password        string `json:"password" validate:"required,min=8,max=128,secure_password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=128,secure_password"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type GoogleOAuthRequest struct {
	Token string `json:"token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

func ValidateSecurePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	result := validators.ValidatePassword(password)
	return result.IsValid
}

func RegisterPasswordValidators(v *validator.Validate) {
	v.RegisterValidation("secure_password", ValidateSecurePassword)
}
