package giftcard

type CreateGiftcardRequest struct {
	RedeemCode string  `json:"redeem_code" validate:"omitempty,min=4,max=50"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

type ApplyGiftcardRequest struct {
	RedeemCode string `json:"redeem_code" validate:"required"`
}

type VerifyGiftcardRequest struct {
	RedeemCode string `json:"redeem_code" validate:"required"`
}

type GiftcardResponse struct {
	ID         string  `json:"id"`
	RedeemCode string  `json:"redeem_code"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	UserEmail  string  `json:"user_email,omitempty"`
	UsedAt     string  `json:"used_at,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func ToResponse(g *Giftcard) *GiftcardResponse {
	usedAt := ""
	if g.UsedAt != nil {
		usedAt = g.UsedAt.Format("2006-01-02T15:04:05Z")
	}

	return &GiftcardResponse{
		ID:         g.ID,
		RedeemCode: g.RedeemCode,
		Amount:     g.Amount,
		Status:     g.Status,
		UserEmail:  g.UserEmail,
		UsedAt:     usedAt,
		CreatedAt:  g.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  g.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(giftcards []*Giftcard) []*GiftcardResponse {
	responses := make([]*GiftcardResponse, len(giftcards))
	for i, g := range giftcards {
		responses[i] = ToResponse(g)
	}
	return responses
}
