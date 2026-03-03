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
	Name             string  `json:"name" validate:"omitempty,min=2,max=100"`
	Email            string  `json:"email" validate:"omitempty,email"`
	Role             *string `json:"role" validate:"omitempty,oneof=user admin"`
	IsActive         *bool   `json:"is_active" validate:"omitempty"`
	ProfilePhoto     *string `json:"profile_photo" validate:"omitempty,url"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty"`
	TelegramLink     *string `json:"telegram_link" validate:"omitempty"`
}

type UserResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	IsActive         bool   `json:"is_active"`
	ProfilePhoto     string `json:"profile_photo"`
	TelegramUsername string `json:"telegram_username"`
	TelegramLink     string `json:"telegram_link"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

func ToResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		Role:             u.Role,
		IsActive:         u.IsActive,
		ProfilePhoto:     u.ProfilePhoto,
		TelegramUsername: u.TelegramUsername,
		TelegramLink:     u.TelegramLink,
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(users []*User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToResponse(u)
	}
	return responses
}
