package wallet

type CreateWalletRequest struct {
	CoinName      string `json:"coin_name" validate:"required,min=2,max=50"`
	Network       string `json:"network" validate:"required,min=2,max=50"`
	CoinLogoURL   string `json:"coin_logo_url" validate:"required"`
	WalletAddress string `json:"wallet_address" validate:"required,min=5,max=255"`
	QrCodeURL     string `json:"qr_code_url" validate:"required"`
}

type UpdateWalletRequest struct {
	CoinName      *string `json:"coin_name" validate:"omitempty,min=2,max=50"`
	Network       *string `json:"network" validate:"omitempty,min=2,max=50"`
	CoinLogoURL   *string `json:"coin_logo_url" validate:"omitempty"`
	WalletAddress *string `json:"wallet_address" validate:"omitempty,min=5,max=255"`
	QrCodeURL     *string `json:"qr_code_url" validate:"omitempty"`
}

type WalletResponse struct {
	ID            string `json:"id"`
	CoinName      string `json:"coin_name"`
	Network       string `json:"network"`
	CoinLogoURL   string `json:"coin_logo_url"`
	WalletAddress string `json:"wallet_address"`
	QrCodeURL     string `json:"qr_code_url"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func ToResponse(w *Wallet) *WalletResponse {
	return &WalletResponse{
		ID:            w.ID,
		CoinName:      w.CoinName,
		Network:       w.Network,
		CoinLogoURL:   w.CoinLogoURL,
		WalletAddress: w.WalletAddress,
		QrCodeURL:     w.QrCodeURL,
		CreatedAt:     w.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     w.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(wallets []*Wallet) []*WalletResponse {
	responses := make([]*WalletResponse, len(wallets))
	for i, w := range wallets {
		responses[i] = ToResponse(w)
	}
	return responses
}
