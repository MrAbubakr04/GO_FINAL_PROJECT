package transaction

import (
	"context"
	"testing"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/repository"
	"github.com/jackc/pgx/v5"
)

type stubTransactionRepo struct{}

func (s *stubTransactionRepo) GetPhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error) {
	return nil, nil
}
func (s *stubTransactionRepo) GetAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
	return nil, nil
}
func (s *stubTransactionRepo) GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error) {
	return nil, nil
}
func (s *stubTransactionRepo) GetAccountForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*domain.Account, error) {
	return nil, nil
}
func (s *stubTransactionRepo) GetClientByID(ctx context.Context, clientID int64) (*domain.Client, error) {
	return nil, nil
}
func (s *stubTransactionRepo) IsAccountIdentified(ctx context.Context, tx pgx.Tx, accountID int64) (bool, error) {
	return false, nil
}
func (s *stubTransactionRepo) UpdateAccountBalance(ctx context.Context, tx pgx.Tx, accountID int64, delta int64) error {
	return nil
}
func (s *stubTransactionRepo) CreateTransaction(ctx context.Context, tx pgx.Tx, txn domain.Transaction) (*domain.Transaction, error) {
	return &txn, nil
}
func (s *stubTransactionRepo) BeginTx(ctx context.Context) (pgx.Tx, error) { return nil, nil }
func (s *stubTransactionRepo) ListAccountTransactions(ctx context.Context, accountID int64) ([]domain.Transaction, error) {
	return nil, nil
}
func (s *stubTransactionRepo) ListClientTransactions(ctx context.Context, clientID int64) ([]domain.Transaction, error) {
	return nil, nil
}

var _ repository.TransactionRepository = (*stubTransactionRepo)(nil)

func TestDepositRejectsInvalidAmount(t *testing.T) {
	svc := NewService(&stubTransactionRepo{})
	_, err := svc.Deposit(context.Background(), DepositInput{Phone: "+998901234567", Amount: 0})
	if err == nil {
		t.Fatal("expected invalid amount error")
	}
	if err != domain.ErrInvalidAmount {
		t.Fatalf("expected %v, got %v", domain.ErrInvalidAmount, err)
	}
}

func TestTransferRejectsSameAccount(t *testing.T) {
	svc := NewService(&stubTransactionRepo{})
	_, err := svc.Transfer(context.Background(), TransferInput{FromPhone: "+998901234567", ToPhone: "+998901234567", Amount: 10})
	if err == nil {
		t.Fatal("expected same account error")
	}
	if err != domain.ErrSameAccountTransfer {
		t.Fatalf("expected %v, got %v", domain.ErrSameAccountTransfer, err)
	}
}
