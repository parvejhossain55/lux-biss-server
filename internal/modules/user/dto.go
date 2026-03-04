package user

type CreateUserRequest struct {
	Name             string `json:"name" validate:"required,min=2,max=100"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8,max=128"`
	Role             string `json:"role" validate:"omitempty,oneof=user admin"`
	ProfilePhoto     string `json:"profile_photo" validate:"omitempty,url"`
	TelegramUsername string `json:"telegram_username" validate:"omitempty"`
	TelegramLink     string `json:"telegram_link" validate:"omitempty"`
}

type UpdateUserRequest struct {
	// Basic
	Name             string  `json:"name" validate:"omitempty,min=2,max=100"`
	Email            string  `json:"email" validate:"omitempty,email"`
	Role             *string `json:"role" validate:"omitempty,oneof=user admin"`
	IsActive         *bool   `json:"is_active" validate:"omitempty"`
	ProfilePhoto     *string `json:"profile_photo" validate:"omitempty,url"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty"`
	TelegramLink     *string `json:"telegram_link" validate:"omitempty"`
	// Personal Information
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty"`
	Gender      *string `json:"gender" validate:"omitempty,oneof=Male Female Other"`
	Phone       *string `json:"phone" validate:"omitempty,max=30"`
	Address     *string `json:"address" validate:"omitempty,max=500"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	// Payment Wallet Information
	PaymentMethod     *string `json:"payment_method" validate:"omitempty,max=100"`
	PaymentCurrency   *string `json:"payment_currency" validate:"omitempty,max=50"`
	PaymentNetwork    *string `json:"payment_network" validate:"omitempty,max=100"`
	WithdrawalAddress *string `json:"withdrawal_address" validate:"omitempty,max=255"`
}

type UserResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Role             string  `json:"role"`
	IsActive         bool    `json:"is_active"`
	ProfilePhoto     string  `json:"profile_photo"`
	TelegramUsername string  `json:"telegram_username"`
	TelegramLink     string  `json:"telegram_link"`
	Balance          float64 `json:"balance"`
	// Personal Information
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	// Payment Wallet Information
	PaymentMethod     string `json:"payment_method"`
	PaymentCurrency   string `json:"payment_currency"`
	PaymentNetwork    string `json:"payment_network"`
	WithdrawalAddress string `json:"withdrawal_address"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func ToResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:                u.ID,
		Name:              u.Name,
		Email:             u.Email,
		Role:              u.Role,
		IsActive:          u.IsActive,
		ProfilePhoto:      u.ProfilePhoto,
		TelegramUsername:  u.TelegramUsername,
		TelegramLink:      u.TelegramLink,
		Balance:           u.Balance,
		DateOfBirth:       u.DateOfBirth,
		Gender:            u.Gender,
		Phone:             u.Phone,
		Address:           u.Address,
		Country:           u.Country,
		PaymentMethod:     u.PaymentMethod,
		PaymentCurrency:   u.PaymentCurrency,
		PaymentNetwork:    u.PaymentNetwork,
		WithdrawalAddress: u.WithdrawalAddress,
		CreatedAt:         u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(users []*User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToResponse(u)
	}
	return responses
}
