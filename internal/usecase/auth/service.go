package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type TokenService interface {
	GeneratePair(ctx context.Context, phoneNum string) (*domain.TokenPair, error)
}

type Service struct {
	repo         repository.AuthRepository
	tokenService TokenService
}

type CheckPhoneResult struct {
	Exists bool   `json:"exists"`
	Status string `json:"status,omitempty"`
}

type LoginInput struct {
	Phone string `json:"phone"`
	PIN   string `json:"pin"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterInput struct {
	Phone  string `json:"phone"`
	PIN    string `json:"pin"`
	Device string `json:"device"`
}

var phonePattern = regexp.MustCompile(`^\+?[0-9]{7,15}$`)

func NewService(repo repository.AuthRepository, tokenService TokenService) *Service {
	if tokenService == nil {
		tokenService = NewJWTService()
	}
	return &Service{repo: repo, tokenService: tokenService}
}

func (s *Service) CheckPhone(ctx context.Context, phoneNum string) (CheckPhoneResult, error) {
	phoneNum = normalizePhone(phoneNum)
	logger.Info("auth check phone started", logger.Fields{"phone": phoneNum})

	phone, err := s.repo.GetActivePhoneByNumber(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to find active phone", logger.Fields{"phone": phoneNum}, err)
		return CheckPhoneResult{}, err
	}
	if phone == nil {
		logger.Info("phone not found", logger.Fields{"phone": phoneNum})
		return CheckPhoneResult{Exists: false}, nil
	}

	account, err := s.repo.GetActiveAccountByPhone(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to find active account", logger.Fields{"phone": phoneNum}, err)
		return CheckPhoneResult{}, err
	}

	status := ""
	if account != nil && account.StatusCode != "" {
		status = strings.ToUpper(account.StatusCode)
	}
	logger.Info("phone check completed", logger.Fields{"phone": phoneNum, "status": status})
	return CheckPhoneResult{Exists: true, Status: status}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*domain.TokenPair, error) {
	phoneNum := normalizePhone(input.Phone)
	logger.Info("auth login started", logger.Fields{"phone": phoneNum})

	if err := validateAuthInput(phoneNum, input.PIN); err != nil {
		logger.Warning("invalid login payload", logger.Fields{"phone": phoneNum})
		return nil, err
	}

	phone, err := s.repo.GetActivePhoneByNumber(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to get phone for login", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	if phone == nil {
		logger.Warning("phone not found during login", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrPhoneNotFound
	}

	account, err := s.repo.GetActiveAccountByPhone(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to get account for login", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	if account == nil {
		logger.Warning("account not found during login", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrPhoneNotFound
	}
	if account.StatusCode != "active" {
		logger.Warning("login blocked because account status is not active", logger.Fields{"phone": phoneNum, "status": account.StatusCode})
		return nil, domain.ErrPhoneInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PINHash), []byte(input.PIN)); err != nil {
		logger.Warning("invalid pin supplied", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrInvalidPin
	}

	tokens, err := s.tokenService.GeneratePair(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to generate tokens", logger.Fields{"phone": phoneNum}, err)
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	logger.Info("login completed", logger.Fields{"phone": phoneNum})
	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*domain.TokenPair, error) {
	logger.Info("auth refresh started", logger.Fields{"refresh_token": strings.TrimSpace(input.RefreshToken)})
	if strings.TrimSpace(input.RefreshToken) == "" {
		return nil, domain.ErrInvalidInput
	}
	if s.tokenService == nil {
		return nil, domain.ErrInvalidInput
	}
	pair, err := s.tokenService.GeneratePair(ctx, "refresh")
	if err != nil {
		logger.Error("failed to generate refresh tokens", logger.Fields{"refresh_token": input.RefreshToken}, err)
		return nil, err
	}
	logger.Info("auth refresh completed", logger.Fields{})
	return pair, nil
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*domain.TokenPair, error) {
	phoneNum := normalizePhone(input.Phone)
	logger.Info("auth registration started", logger.Fields{"phone": phoneNum})

	if err := validateAuthInput(phoneNum, input.PIN); err != nil {
		logger.Warning("invalid registration payload", logger.Fields{"phone": phoneNum})
		return nil, err
	}
	if strings.TrimSpace(input.Device) == "" {
		logger.Warning("missing device on registration", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrInvalidInput
	}

	phone, err := s.repo.GetActivePhoneByNumber(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to verify phone before registration", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	if phone != nil {
		logger.Warning("registration blocked because phone already exists", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrPhoneAlreadyExists
	}

	account, err := s.repo.GetActiveAccountByPhone(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to verify account before registration", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	if account != nil {
		logger.Warning("registration blocked because account already exists", logger.Fields{"phone": phoneNum})
		return nil, domain.ErrAccountAlreadyExists
	}

	pinHash, err := bcrypt.GenerateFromPassword([]byte(input.PIN), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash pin", logger.Fields{"phone": phoneNum}, err)
		return nil, fmt.Errorf("hash pin: %w", err)
	}

	statusID, err := s.repo.GetStatusIDByCode(ctx, "active")
	if err != nil {
		logger.Error("failed to get active status id", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	if statusID == 0 {
		logger.Error("active status id not found", logger.Fields{"phone": phoneNum}, nil)
		return nil, fmt.Errorf("active status id not found")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed to start registration transaction", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	createdPhone, err := s.repo.CreatePhone(ctx, tx, domain.Phone{PhoneNum: phoneNum})
	if err != nil {
		logger.Error("failed to create phone during registration", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}

	_, err = s.repo.CreateAccount(ctx, tx, domain.Account{
		PhoneNum:  phoneNum,
		PINHash:   string(pinHash),
		BalanceTJ: 0,
		BalanceRU: 0,
		BalanceEN: 0,
		Device:    input.Device,
		IsActive:  true,
		StatusID:  statusID,
	})
	if err != nil {
		logger.Error("failed to create account during registration", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit registration transaction", logger.Fields{"phone": phoneNum}, err)
		return nil, err
	}

	tokens, err := s.tokenService.GeneratePair(ctx, phoneNum)
	if err != nil {
		logger.Error("failed to generate tokens after registration", logger.Fields{"phone": phoneNum}, err)
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	logger.Info("registration completed", logger.Fields{"phone": phoneNum, "phone_id": createdPhone.ID, "created_at": time.Now().UTC().Format(time.RFC3339)})
	return tokens, nil
}

func normalizePhone(phone string) string {
	return strings.TrimSpace(strings.ReplaceAll(phone, " ", ""))
}

func validateAuthInput(phone, pin string) error {
	if !phonePattern.MatchString(phone) {
		return domain.ErrInvalidInput
	}
	if len(strings.TrimSpace(pin)) < 4 || len(strings.TrimSpace(pin)) > 32 {
		return domain.ErrInvalidInput
	}
	return nil
}
