package user

type ManagerResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	TelegramUsername string `json:"telegram_username"`
	ProfilePhoto     string `json:"profile_photo"`
}

type LevelResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StepResponse struct {
	ID         uint   `json:"id"`
	LevelID    uint   `json:"level_id"`
	StepNumber int    `json:"step_number"`
	Name       string `json:"name"`
}

type CreateUserRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8,max=128"`
	Role         string `json:"role" validate:"omitempty,oneof=user admin"`
	ProfilePhoto string `json:"profile_photo" validate:"omitempty"`
}

type UpdateUserRequest struct {
	// Basic
	Name                string   `json:"name" validate:"omitempty,min=2,max=100"`
	Email               string   `json:"email" validate:"omitempty,email"`
	Role                *string  `json:"role" validate:"omitempty,oneof=user admin"`
	Status              *string  `json:"status" validate:"omitempty,oneof=active ignored suspend hold"`
	ProfilePhoto        *string  `json:"profile_photo" validate:"omitempty"`
	Balance             *float64 `json:"balance" validate:"omitempty"`
	HoldBalance         *float64 `json:"hold_balance" validate:"omitempty"`
	WithdrawableBalance *float64 `json:"withdrawable_balance" validate:"omitempty"`
	// Personal Information
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty"`
	Gender      *string `json:"gender" validate:"omitempty,oneof=Male Female Other"`
	Phone       *string `json:"phone" validate:"omitempty,max=30"`
	Address     *string `json:"address" validate:"omitempty,max=500"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	// Payment Wallet Information
	PaymentMethod        *string `json:"payment_method" validate:"omitempty,max=100"`
	PaymentCurrency      *string `json:"payment_currency" validate:"omitempty,max=50"`
	PaymentNetwork       *string `json:"payment_network" validate:"omitempty,max=100"`
	WithdrawalAddress    *string `json:"withdrawal_address" validate:"omitempty,max=255"`
	LevelID              *uint   `json:"level_id" validate:"omitempty"`
	StepID               *uint   `json:"step_id" validate:"omitempty"`
	CurrentStepCompleted *bool   `json:"current_step_completed" validate:"omitempty"`
}

type UserResponse struct {
	ID                   string           `json:"id"`
	Name                 string           `json:"name"`
	Email                string           `json:"email"`
	Role                 string           `json:"role"`
	Status               string           `json:"status"`
	ProfilePhoto         string           `json:"profile_photo"`
	Balance              float64          `json:"balance"`
	HoldBalance          float64          `json:"hold_balance"`
	WithdrawableBalance  float64          `json:"withdrawable_balance"`
	ManagerID            *string          `json:"manager_id"`
	Manager              *ManagerResponse `json:"manager,omitempty"`
	LevelID              *uint            `json:"level_id"`
	Level                *LevelResponse   `json:"level,omitempty"`
	StepID               *uint            `json:"step_id"`
	Step                 *StepResponse    `json:"step,omitempty"`
	CurrentStepCompleted bool             `json:"current_step_completed"`
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
	var managerResponse *ManagerResponse
	if u.Manager != nil {
		managerResponse = &ManagerResponse{
			ID:               u.Manager.ID,
			Name:             u.Manager.Name,
			TelegramUsername: u.Manager.TelegramUsername,
			ProfilePhoto:     u.Manager.ProfilePhoto,
		}
	}

	var levelResponse *LevelResponse
	if u.Level != nil {
		levelResponse = &LevelResponse{
			ID:   u.Level.ID,
			Name: u.Level.Name,
		}
	}

	var stepResponse *StepResponse
	if u.Step != nil {
		stepResponse = &StepResponse{
			ID:         u.Step.ID,
			LevelID:    u.Step.LevelID,
			StepNumber: u.Step.StepNumber,
			Name:       u.Step.Name,
		}
	}

	return &UserResponse{
		ID:                   u.ID,
		Name:                 u.Name,
		Email:                u.Email,
		Role:                 u.Role,
		Status:               u.Status,
		ProfilePhoto:         u.ProfilePhoto,
		Balance:              u.Balance,
		HoldBalance:          u.HoldBalance,
		WithdrawableBalance:  u.WithdrawableBalance,
		ManagerID:            u.ManagerID,
		Manager:              managerResponse,
		LevelID:              u.LevelID,
		Level:                levelResponse,
		StepID:               u.StepID,
		Step:                 stepResponse,
		CurrentStepCompleted: u.CurrentStepCompleted,
		DateOfBirth:          u.DateOfBirth,
		Gender:               u.Gender,
		Phone:                u.Phone,
		Address:              u.Address,
		Country:              u.Country,
		PaymentMethod:        u.PaymentMethod,
		PaymentCurrency:      u.PaymentCurrency,
		PaymentNetwork:       u.PaymentNetwork,
		WithdrawalAddress:    u.WithdrawalAddress,
		CreatedAt:            u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:            u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(users []*User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToResponse(u)
	}
	return responses
}
