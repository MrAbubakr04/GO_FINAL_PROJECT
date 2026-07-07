package auth

import (
	"context"
	"testing"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/jackc/pgx/v5"
)

type stubRepo struct {
	phone *domain.Phone
	acct  *domain.Account
}

func (s *stubRepo) GetActivePhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error) {
	return s.phone, nil
}

func (s *stubRepo) GetActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
	return s.acct, nil
}

func (s *stubRepo) GetStatusIDByCode(ctx context.Context, code string) (int, error) {
	return 1, nil
}

func (s *stubRepo) GetStatusCodeByID(ctx context.Context, id int) (string, error) {
	return "active", nil
}

func (s *stubRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (s *stubRepo) CreatePhone(ctx context.Context, tx pgx.Tx, phone domain.Phone) (*domain.Phone, error) {
	return &phone, nil
}

func (s *stubRepo) CreateAccount(ctx context.Context, tx pgx.Tx, account domain.Account) (*domain.Account, error) {
	return &account, nil
}

func TestCheckPhoneReturnsExistsAndStatus(t *testing.T) {
	repo := &stubRepo{phone: &domain.Phone{PhoneNum: "+998901234567"}, acct: &domain.Account{StatusCode: "active"}}
	service := NewService(repo, nil)

	result, err := service.CheckPhone(context.Background(), "+998901234567")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !result.Exists {
		t.Fatal("expected exists true")
	}
	if result.Status != "ACTIVE" {
		t.Fatalf("expected status ACTIVE, got %s", result.Status)
	}
}

func TestCheckPhoneReturnsFalseWhenPhoneMissing(t *testing.T) {
	repo := &stubRepo{}
	service := NewService(repo, nil)

	result, err := service.CheckPhone(context.Background(), "+998901234567")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Exists {
		t.Fatal("expected exists false")
	}
}

func TestAuthorizeFailsForInvalidPIN(t *testing.T) {
	repo := &stubRepo{phone: &domain.Phone{PhoneNum: "+998901234567"}, acct: &domain.Account{PhoneNum: "+998901234567", PINHash: "$2a$10$abcdefghijklmnopqrstuv"}}
	service := NewService(repo, nil)

	_, err := service.Login(context.Background(), LoginInput{Phone: "+998901234567", PIN: "1234"})
	if err == nil {
		t.Fatal("expected invalid pin error")
	}
}

func TestRefreshRejectsEmptyToken(t *testing.T) {
	service := NewService(&stubRepo{}, nil)

	_, err := service.Refresh(context.Background(), RefreshInput{})
	if err == nil {
		t.Fatal("expected invalid input error")
	}
}
