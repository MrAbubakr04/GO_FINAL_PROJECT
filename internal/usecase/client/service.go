package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.ClientRepository
}

type LoginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateClientInput struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Fathername string `json:"fathername"`
	DocNum     string `json:"doc_num"`
	TIN        string `json:"tin"`
	BirthDate  string `json:"birth_date"`
	Gender     string `json:"gender"`
	Address    string `json:"address"`
}

type UpdateClientInput struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Fathername string `json:"fathername"`
	DocNum     string `json:"doc_num"`
	TIN        string `json:"tin"`
	BirthDate  string `json:"birth_date"`
	Gender     string `json:"gender"`
	Address    string `json:"address"`
}

type IdentifyInput struct {
	Phone string `json:"phone"`
	TIN   string `json:"tin"`
}

func NewService(repo repository.ClientRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*domain.User, error) {
	logger.Info("employee login started", logger.Fields{"login": input.Login})
	user, err := s.repo.GetUserByLogin(ctx, input.Login)
	if err != nil {
		logger.Error("failed to fetch employee", logger.Fields{"login": input.Login}, err)
		return nil, err
	}
	if user == nil {
		logger.Warning("employee not found", logger.Fields{"login": input.Login})
		return nil, domain.ErrUserNotFound
	}
	if user.ActiveTo != nil {
		logger.Warning("employee inactive", logger.Fields{"login": input.Login})
		return nil, domain.ErrUserInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		logger.Warning("invalid password", logger.Fields{"login": input.Login})
		return nil, domain.ErrInvalidPassword
	}
	logger.Info("employee login completed", logger.Fields{"login": input.Login})
	return user, nil
}

func (s *Service) CreateClient(ctx context.Context, input CreateClientInput) (*domain.Client, error) {
	logger.Info("client registration started", logger.Fields{"tin": input.TIN})
	if strings.TrimSpace(input.TIN) == "" {
		return nil, domain.ErrInvalidInput
	}
	if existing, err := s.repo.GetClientByTIN(ctx, input.TIN); err != nil {
		logger.Error("failed to check client tin", logger.Fields{"tin": input.TIN}, err)
		return nil, err
	} else if existing != nil {
		logger.Warning("client already exists", logger.Fields{"tin": input.TIN})
		return nil, domain.ErrClientAlreadyExists
	}
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		logger.Warning("invalid birth date", logger.Fields{"tin": input.TIN})
		return nil, domain.ErrInvalidInput
	}
	client := domain.Client{
		Name:       strings.TrimSpace(input.Name),
		Surname:    strings.TrimSpace(input.Surname),
		Fathername: strings.TrimSpace(input.Fathername),
		DocNum:     strings.TrimSpace(input.DocNum),
		TIN:        strings.TrimSpace(input.TIN),
		BirthDate:  birthDate,
		Gender:     strings.TrimSpace(input.Gender),
		Address:    strings.TrimSpace(input.Address),
	}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	created, err := s.repo.CreateClient(ctx, tx, client)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	logger.Info("client registration completed", logger.Fields{"client_id": created.ID})
	return created, nil
}

func (s *Service) GetClientByID(ctx context.Context, id int64) (*domain.Client, error) {
	client, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, domain.ErrClientNotFound
	}
	return client, nil
}

func (s *Service) GetClientByTIN(ctx context.Context, tin string) (*domain.Client, error) {
	client, err := s.repo.GetClientByTIN(ctx, tin)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, domain.ErrClientNotFound
	}
	return client, nil
}

func (s *Service) UpdateClient(ctx context.Context, id int64, input UpdateClientInput) (*domain.Client, error) {
	current, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, domain.ErrClientNotFound
	}
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}
	updated := domain.Client{ID: id, Name: strings.TrimSpace(input.Name), Surname: strings.TrimSpace(input.Surname), Fathername: strings.TrimSpace(input.Fathername), DocNum: strings.TrimSpace(input.DocNum), TIN: strings.TrimSpace(input.TIN), BirthDate: birthDate, Gender: strings.TrimSpace(input.Gender), Address: strings.TrimSpace(input.Address)}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	result, err := s.repo.UpdateClient(ctx, tx, updated)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, domain.ErrClientInactive
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DeactivateClient(ctx context.Context, clientID int64) error {
	logger.Info("client deactivation started", logger.Fields{"client_id": clientID})
	client, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return err
	}
	if client == nil {
		return domain.ErrClientNotFound
	}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	if err = s.repo.DeactivateClient(ctx, tx, clientID); err != nil {
		return err
	}
	accounts, err := s.repo.GetActiveAccountsByClientID(ctx, clientID)
	if err != nil {
		return err
	}
	for _, acc := range accounts {
		if err = s.repo.DeactivateAccount(ctx, tx, acc.ID); err != nil {
			return err
		}
	}
	phone, err := s.repo.GetActivePhoneByNumber(ctx, "")
	_ = phone
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	logger.Info("client deactivation completed", logger.Fields{"client_id": clientID})
	return nil
}

func (s *Service) IdentifyClient(ctx context.Context, input IdentifyInput) (bool, error) {
	phoneNum := strings.TrimSpace(input.Phone)
	tin := strings.TrimSpace(input.TIN)
	if phoneNum == "" || tin == "" {
		return false, domain.ErrInvalidInput
	}

	client, err := s.repo.GetClientByTIN(ctx, tin)
	if err != nil {
		return false, err
	}
	if client == nil {
		return false, domain.ErrClientNotFound
	}
	phone, err := s.repo.GetActivePhoneByNumber(ctx, phoneNum)
	if err != nil {
		return false, err
	}
	if phone == nil {
		return false, domain.ErrPhoneNotFound
	}
	if phone.ClientID != nil && *phone.ClientID != client.ID {
		return false, domain.ErrPhoneAlreadyUsed
	}
	account, err := s.repo.GetActiveAccountByPhone(ctx, phoneNum)
	if err != nil {
		return false, err
	}
	if account == nil {
		return false, domain.ErrAccountNotFound
	}
	if phone.ClientID == nil {
		tx, err := s.repo.BeginTx(ctx)
		if err != nil {
			return false, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback(ctx)
			}
		}()

		if err = s.repo.AttachPhoneToClient(ctx, tx, phone.ID, client.ID); err != nil {
			return false, err
		}
		if err = tx.Commit(ctx); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s *Service) CreateEmployee(ctx context.Context, login, password string) (*domain.User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u := domain.User{Login: login, PasswordHash: string(pwdHash), Role: "employee"}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	created, err := s.repo.CreateUser(ctx, tx, u)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}
