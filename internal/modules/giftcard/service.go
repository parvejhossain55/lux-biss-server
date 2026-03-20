package giftcard

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/modules/transaction"
	"github.com/parvej/luxbiss_server/internal/modules/user"
)

type GiftcardService struct {
	repo           Repository
	userService    user.Service
	txRepo         transaction.Repository
	log            *logger.Logger
	telegramToken  string
	telegramChatID string
	telegramProxy  string
}

func NewService(repo Repository, userService user.Service, txRepo transaction.Repository, log *logger.Logger, token, chatID, proxy string) *GiftcardService {
	return &GiftcardService{
		repo:           repo,
		userService:    userService,
		txRepo:         txRepo,
		log:            log,
		telegramToken:  token,
		telegramChatID: chatID,
		telegramProxy:  proxy,
	}
}

func GenerateRedeemCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return fmt.Sprintf("%s-%s-%s-%s", string(b[0:4]), string(b[4:8]), string(b[8:12]), string(b[12:16]))
}

func (s *GiftcardService) Create(ctx context.Context, req *CreateGiftcardRequest) (*Giftcard, error) {
	code := strings.ToUpper(strings.TrimSpace(req.RedeemCode))
	if code == "" {
		code = GenerateRedeemCode()
	}

	giftcard := &Giftcard{
		RedeemCode: code,
		Amount:     req.Amount,
		Status:     StatusAvailable,
	}

	if req.ExpiredAt != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", req.ExpiredAt)
		if err == nil {
			giftcard.ExpiredAt = &t
		}
	}

	if err := s.repo.Create(ctx, giftcard); err != nil {
		s.log.Errorw("Failed to create giftcard", "error", err, "code", code)
		return nil, common.ErrInternal(err)
	}

	s.log.Infow("Giftcard created successfully", "giftcard_id", giftcard.ID)
	return giftcard, nil
}

func (s *GiftcardService) List(ctx context.Context, limit, offset int) ([]*Giftcard, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *GiftcardService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Errorw("Failed to delete giftcard", "error", err, "giftcard_id", id)
		return err
	}
	s.log.Infow("Giftcard deleted successfully", "giftcard_id", id)
	return nil
}

func (s *GiftcardService) Apply(ctx context.Context, req *ApplyGiftcardRequest, userID, userEmail string) (*Giftcard, error) {
	code := strings.ToUpper(strings.TrimSpace(req.RedeemCode))
	giftcard, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, common.ErrNotFound("Giftcard")
	}

	if giftcard.Status != StatusAvailable {
		return nil, common.ErrBadRequest("Giftcard has already been used")
	}

	if giftcard.ExpiredAt != nil && giftcard.ExpiredAt.Before(time.Now()) {
		return nil, common.ErrBadRequest("Giftcard has expired")
	}

	now := time.Now()

	giftcard.Status = StatusUsed
	giftcard.UserID = &userID
	giftcard.UserEmail = userEmail
	giftcard.UsedAt = &now

	if err := s.repo.Update(ctx, giftcard); err != nil {
		s.log.Errorw("Failed to apply giftcard", "error", err, "giftcard_id", giftcard.ID)
		return nil, common.ErrInternal(err)
	}

	// Create a pending deposit transaction. The balance will be updated when an admin approves it.
	tx := &transaction.Transaction{
		UserID: userID,
		Type:   transaction.TypeDeposit,
		Amount: giftcard.Amount,
		Status: transaction.StatusCompleted,
		TxHash: common.GenerateHash(),
		Note:   "Gift card redeemed: " + giftcard.RedeemCode,
	}
	if err := s.txRepo.Create(ctx, tx); err != nil {
		s.log.Errorw("Giftcard applied but failed to create transaction record", "error", err, "user_id", userID)
		return nil, common.ErrInternal(err)
	}

	// Add to user's hold balance
	if err := s.userService.UpdateHoldBalance(ctx, userID, giftcard.Amount); err != nil {
		s.log.Errorw("Giftcard applied but failed to update user hold balance", "error", err, "user_id", userID)
	}

	// Send Telegram notification
	go func() {
		msg := fmt.Sprintf("🎁 <b>Gift Card Redeemed</b>\n\n👤 <b>User:</b> %s\n🏷 <b>Code:</b> <code>%s</code>\n💰 <b>Amount:</b> $%s\n📅 <b>Time:</b> %s",
			userEmail,
			giftcard.RedeemCode,
			fmt.Sprintf("%.2f", giftcard.Amount),
			time.Now().Format("2006-01-02 03:04:05 PM"),
		)
		if err := common.SendTelegramMessage(s.telegramToken, s.telegramChatID, msg, s.telegramProxy); err != nil {
			s.log.Errorw("Failed to send telegram notification for giftcard redemption", "error", err)
		}
	}()

	s.log.Infow("Giftcard applied successfully and added to hold balance", "giftcard_id", giftcard.ID, "user_id", userID, "amount", giftcard.Amount)
	return giftcard, nil
}

func (s *GiftcardService) Verify(ctx context.Context, req *VerifyGiftcardRequest) (*Giftcard, error) {
	code := strings.ToUpper(strings.TrimSpace(req.RedeemCode))
	giftcard, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, common.ErrNotFound("Giftcard")
	}

	if giftcard.Status != StatusAvailable {
		return nil, common.ErrBadRequest("Giftcard has already been used")
	}

	if giftcard.ExpiredAt != nil && giftcard.ExpiredAt.Before(time.Now()) {
		return nil, common.ErrBadRequest("Giftcard has expired")
	}

	return giftcard, nil
}
