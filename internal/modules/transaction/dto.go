package transaction

type CreateTransactionRequest struct {
	UserID string  `json:"user_id" validate:"omitempty,uuid4"`
	Type   string  `json:"type" validate:"required,oneof=deposit withdraw"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	TxHash string  `json:"tx_hash" validate:"omitempty,max=255"`
	Note   string  `json:"note" validate:"omitempty,max=1000"`
}

type UpdateTransactionRequest struct {
	Type   string  `json:"type" validate:"omitempty,oneof=deposit withdraw"`
	Amount float64 `json:"amount" validate:"omitempty,gt=0"`
	Status string  `json:"status" validate:"omitempty,oneof=processing completed rejected cancelled"`
	TxHash string  `json:"tx_hash" validate:"omitempty,max=255"`
	Note   string  `json:"note" validate:"omitempty,max=1000"`
}

type InvestRequest struct {
	LevelID  uint `json:"level_id" validate:"required,gt=0"`
	StepID   uint `json:"step_id" validate:"required,gt=0"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}

type TransactionResponse struct {
	ID string `json:"id"`

	UserID    string  `json:"user_id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	TxHash    string  `json:"tx_hash"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type SummaryResponse struct {
	AvailableBalance float64      `json:"available_balance"`
	HoldBalance      float64      `json:"hold_balance"`
	TotalDeposit     float64      `json:"total_deposit"`
	TotalWithdrawal  float64      `json:"total_withdrawal"`
	PeriodDays       int          `json:"period_days"`
	WithdrawReport   []ReportItem `json:"withdraw_report"`
}

func ToResponse(t *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID: t.ID,

		UserID:    t.UserID,
		Type:      t.Type,
		Amount:    t.Amount,
		Status:    t.Status,
		TxHash:    t.TxHash,
		Note:      t.Note,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToResponseList(txs []*Transaction) []*TransactionResponse {
	responses := make([]*TransactionResponse, len(txs))
	for i, t := range txs {
		responses[i] = ToResponse(t)
	}
	return responses
}
