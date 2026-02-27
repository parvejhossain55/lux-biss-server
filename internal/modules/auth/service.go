package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/config"
	"github.com/parvej/luxbiss_server/internal/logger"
	"github.com/parvej/luxbiss_server/internal/modules/user"
	"github.com/parvej/luxbiss_server/pkg/email"
	"github.com/parvej/luxbiss_server/pkg/hash"
	"github.com/parvej/luxbiss_server/pkg/jwt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/idtoken"
)

type Service struct {
	userService user.Service
	jwtManager  *jwt.Manager
	rdb         *redis.Client
	emailSender email.Sender
	oauthCfg    *config.OAuthConfig
	log         *logger.Logger
}

func NewService(
	userService user.Service,
	jwtManager *jwt.Manager,
	rdb *redis.Client,
	emailSender email.Sender,
	oauthCfg *config.OAuthConfig,
	log *logger.Logger,
) *Service {
	return &Service{
		userService: userService,
		jwtManager:  jwtManager,
		rdb:         rdb,
		emailSender: emailSender,
		oauthCfg:    oauthCfg,
		log:         log,
	}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	createReq := &user.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     user.RoleUser,
	}

	newUser, err := s.userService.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokenPair(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		s.log.Errorw("Failed to generate tokens", "error", err, "user_id", newUser.ID)
		return nil, common.ErrInternal(err)
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user.ToResponse(newUser),
	}, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	existingUser, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, common.ErrUnauthorized("Invalid credentials")
	}

	if err := hash.CheckPassword(req.Password, existingUser.Password); err != nil {
		return nil, common.ErrUnauthorized("Invalid credentials")
	}

	if !existingUser.IsActive {
		return nil, common.ErrForbidden("Your account has been deactivated")
	}

	tokens, err := s.jwtManager.GenerateTokenPair(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		s.log.Errorw("Failed to generate tokens", "error", err, "user_id", existingUser.ID)
		return nil, common.ErrInternal(err)
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user.ToResponse(existingUser),
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error) {
	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, common.ErrUnauthorized("Invalid or expired refresh token")
	}

	existingUser, err := s.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, common.ErrUnauthorized("User not found")
	}

	if !existingUser.IsActive {
		return nil, common.ErrForbidden("Your account has been deactivated")
	}

	tokens, err := s.jwtManager.GenerateTokenPair(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		s.log.Errorw("Failed to generate tokens", "error", err, "user_id", existingUser.ID)
		return nil, common.ErrInternal(err)
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user.ToResponse(existingUser),
	}, nil
}

func (s *Service) GoogleLogin(ctx context.Context, req *GoogleOAuthRequest) (*AuthResponse, error) {
	payload, err := idtoken.Validate(ctx, req.Token, s.oauthCfg.GoogleClientID)
	if err != nil {
		s.log.Errorw("Google token validation failed", "error", err)
		return nil, common.ErrUnauthorized("Invalid Google token")
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	existingUser, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		// Register new user
		createReq := &user.CreateUserRequest{
			Name:     name,
			Email:    email,
			Password: generateRandomPassword(), // Google users don't need password but we store a random one
			Role:     user.RoleUser,
		}
		existingUser, err = s.userService.Create(ctx, createReq)
		if err != nil {
			return nil, err
		}
	}

	if !existingUser.IsActive {
		return nil, common.ErrForbidden("Your account has been deactivated")
	}

	tokens, err := s.jwtManager.GenerateTokenPair(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user.ToResponse(existingUser),
	}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	existingUser, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if user exists
		return nil
	}

	otp := generateOTP(6)
	otpKey := fmt.Sprintf("otp:%s", existingUser.Email)

	err = s.rdb.Set(ctx, otpKey, otp, 15*time.Minute).Err()
	if err != nil {
		s.log.Errorw("Failed to store OTP in redis", "error", err)
		return common.ErrInternal(err)
	}

	go func() {
		subject := "Your Password Reset OTP"
		body := fmt.Sprintf("<h1>OTP: %s</h1><p>This code will expire in 15 minutes.</p>", otp)
		err := s.emailSender.SendEmail([]string{existingUser.Email}, subject, body)
		if err != nil {
			s.log.Errorw("Failed to send OTP email", "error", err, "email", existingUser.Email)
		}
	}()

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	otpKey := fmt.Sprintf("otp:%s", req.Email)
	storedOTP, err := s.rdb.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			return common.ErrBadRequest("Invalid or expired OTP")
		}
		return common.ErrInternal(err)
	}

	if storedOTP != req.OTP {
		return common.ErrBadRequest("Invalid OTP")
	}

	existingUser, err := s.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		return common.ErrNotFound("User")
	}

	// Update password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return common.ErrInternal(err)
	}

	// We'll add this method to userService
	if err := s.userService.UpdatePassword(ctx, existingUser.ID, hashedPassword); err != nil {
		return err
	}

	s.rdb.Del(ctx, otpKey)
	return nil
}

func generateOTP(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}
	return string(b)
}

func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	length := 16
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}
	return string(b)
}
