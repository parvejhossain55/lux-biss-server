package manager

type CreateManagerRequest struct {
	Name             string `json:"name" validate:"required,min=2,max=100"`
	TelegramUsername string `json:"telegram_username" validate:"omitempty,max=100"`
	TelegramLink     string `json:"telegram_link" validate:"omitempty"`
	ProfilePhoto     string `json:"profile_photo" validate:"omitempty"`
}

type UpdateManagerRequest struct {
	Name             *string `json:"name" validate:"omitempty,min=2,max=100"`
	TelegramUsername *string `json:"telegram_username" validate:"omitempty,max=100"`
	TelegramLink     *string `json:"telegram_link" validate:"omitempty"`
	ProfilePhoto     *string `json:"profile_photo" validate:"omitempty"`
}

type ManagerResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	TelegramUsername string `json:"telegram_username"`
	TelegramLink     string `json:"telegram_link"`
	ProfilePhoto     string `json:"profile_photo"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

func ToResponse(m *Manager) *ManagerResponse {
	return &ManagerResponse{
		ID:               m.ID,
		Name:             m.Name,
		TelegramUsername: m.TelegramUsername,
		TelegramLink:     m.TelegramLink,
		ProfilePhoto:     m.ProfilePhoto,
		CreatedAt:        m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(managers []*Manager) []*ManagerResponse {
	responses := make([]*ManagerResponse, len(managers))
	for i, m := range managers {
		responses[i] = ToResponse(m)
	}
	return responses
}
