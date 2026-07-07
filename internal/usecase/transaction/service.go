package transaction

import (
	"context"
	"fmt"
	"strings"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/repository"
)

type Service struct {
	repo repository.TransactionRepository
}

type DepositInput struct {
	Phone      string `json:"phone"`
	Amount     int64  `json:"amount"`
	TerminalID string `json:"terminal_id,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Initiator  string `json:"initiator,omitempty"`
}

type TransferInput struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func NewService(repo repository.TransactionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Deposit(ctx context.Context, input DepositInput) (*domain.Transaction, error) {
	logger.Info("transaction deposit started", logger.Fields{"phone": input.Phone, "amount": input.Amount})
	if strings.TrimSpace(input.Phone) == "" || input.Amount <= 0 {
		logger.Warning("invalid deposit input", logger.Fields{"phone": input.Phone, "amount": input.Amount})
		return nil, domain.ErrInvalidAmount
	}
	phone, err := s.repo.GetPhoneByNumber(ctx, input.Phone)
	if err != nil {
		logger.Error("failed to get phone for deposit", logger.Fields{"phone": input.Phone}, err)
		return nil, err
	}
	if phone == nil {
		logger.Warning("phone not found for deposit", logger.Fields{"phone": input.Phone})
		return nil, domain.ErrPhoneNotFound
	}
	if phone.ActiveTo != nil {
		logger.Warning("phone inactive for deposit", logger.Fields{"phone": input.Phone})
		return nil, domain.ErrPhoneInactive
	}
	account, err := s.repo.GetAccountByPhone(ctx, input.Phone)
	if err != nil {
		logger.Error("failed to get account for deposit", logger.Fields{"phone": input.Phone}, err)
		return nil, err
	}
	if account == nil {
		logger.Warning("account not found for deposit", logger.Fields{"phone": input.Phone})
		return nil, domain.ErrAccountNotFound
	}
	if account.ActiveTo != nil || !account.IsActive || account.StatusCode != "active" {
		logger.Warning("account inactive for deposit", logger.Fields{"account_id": account.ID})
		return nil, domain.ErrAccountInactive
	}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed to begin deposit transaction", logger.Fields{"account_id": account.ID}, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	updatedAcc, err := s.repo.GetAccountForUpdate(ctx, tx, account.ID)
	if err != nil {
		logger.Error("failed to lock account for deposit", logger.Fields{"account_id": account.ID}, err)
		return nil, err
	}
	if updatedAcc == nil {
		return nil, domain.ErrAccountNotFound
	}
	if err := s.repo.UpdateAccountBalance(ctx, tx, updatedAcc.ID, input.Amount); err != nil {
		logger.Error("failed to update balance during deposit", logger.Fields{"account_id": updatedAcc.ID}, err)
		return nil, err
	}
	created, err := s.repo.CreateTransaction(ctx, tx, domain.Transaction{FromAccountID: updatedAcc.ID, ToAccountID: updatedAcc.ID, Amount: input.Amount})
	if err != nil {
		logger.Error("failed to create deposit transaction", logger.Fields{"account_id": updatedAcc.ID}, err)
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit deposit transaction", logger.Fields{"account_id": updatedAcc.ID}, err)
		return nil, err
	}
	logger.Info("transaction deposit completed", logger.Fields{"account_id": updatedAcc.ID, "amount": input.Amount})
	return created, nil
}

func (s *Service) Transfer(ctx context.Context, input TransferInput) (*domain.TransferResult, error) {
	logger.Info("transaction transfer started", logger.Fields{"from": input.FromAccountID, "to": input.ToAccountID, "amount": input.Amount})
	if input.FromAccountID == input.ToAccountID {
		logger.Warning("same account transfer rejected", logger.Fields{"from": input.FromAccountID, "to": input.ToAccountID})
		return nil, domain.ErrSameAccountTransfer
	}
	if input.Amount <= 0 {
		logger.Warning("invalid transfer amount", logger.Fields{"amount": input.Amount})
		return nil, domain.ErrInvalidAmount
	}
	if input.FromAccountID == 0 || input.ToAccountID == 0 {
		return nil, domain.ErrInvalidInput
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed to begin transfer transaction", logger.Fields{"from": input.FromAccountID, "to": input.ToAccountID}, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	fromAcc, err := s.repo.GetAccountForUpdate(ctx, tx, input.FromAccountID)
	if err != nil {
		logger.Error("failed to lock sender account", logger.Fields{"account_id": input.FromAccountID}, err)
		return nil, err
	}
	if fromAcc == nil {
		return nil, domain.ErrAccountNotFound
	}
	if fromAcc.ActiveTo != nil || !fromAcc.IsActive || fromAcc.StatusCode != "active" {
		return nil, domain.ErrAccountInactive
	}
	if fromAcc.BalanceTJ < input.Amount {
		return nil, domain.ErrInsufficientFunds
	}
	isIdentified, err := s.repo.IsAccountIdentified(ctx, tx, input.FromAccountID)
	if err != nil {
		logger.Error("failed to verify sender identification", logger.Fields{"account_id": input.FromAccountID}, err)
		return nil, err
	}
	if !isIdentified {
		return nil, domain.ErrAccountNotIdentified
	}
	toAcc, err := s.repo.GetAccountForUpdate(ctx, tx, input.ToAccountID)
	if err != nil {
		logger.Error("failed to lock receiver account", logger.Fields{"account_id": input.ToAccountID}, err)
		return nil, err
	}
	if toAcc == nil {
		return nil, domain.ErrAccountNotFound
	}
	if toAcc.ActiveTo != nil || !toAcc.IsActive || toAcc.StatusCode != "active" {
		return nil, domain.ErrAccountInactive
	}
	toIdentified, err := s.repo.IsAccountIdentified(ctx, tx, input.ToAccountID)
	if err != nil {
		logger.Error("failed to verify receiver identification", logger.Fields{"account_id": input.ToAccountID}, err)
		return nil, err
	}
	if !toIdentified {
		return nil, domain.ErrAccountNotIdentified
	}
	if err := s.repo.UpdateAccountBalance(ctx, tx, fromAcc.ID, -input.Amount); err != nil {
		logger.Error("failed to debit sender account", logger.Fields{"account_id": fromAcc.ID}, err)
		return nil, err
	}
	if err := s.repo.UpdateAccountBalance(ctx, tx, toAcc.ID, input.Amount); err != nil {
		logger.Error("failed to credit receiver account", logger.Fields{"account_id": toAcc.ID}, err)
		return nil, err
	}
	outTxn, err := s.repo.CreateTransaction(ctx, tx, domain.Transaction{FromAccountID: fromAcc.ID, ToAccountID: toAcc.ID, Amount: input.Amount})
	if err != nil {
		logger.Error("failed to create transfer out transaction", logger.Fields{"account_id": fromAcc.ID}, err)
		return nil, err
	}
	inTxn, err := s.repo.CreateTransaction(ctx, tx, domain.Transaction{FromAccountID: fromAcc.ID, ToAccountID: toAcc.ID, Amount: input.Amount})
	if err != nil {
		logger.Error("failed to create transfer in transaction", logger.Fields{"account_id": toAcc.ID}, err)
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit transfer transaction", logger.Fields{"from": input.FromAccountID, "to": input.ToAccountID}, err)
		return nil, err
	}
	logger.Info("transaction transfer completed", logger.Fields{"from": input.FromAccountID, "to": input.ToAccountID, "amount": input.Amount})
	return &domain.TransferResult{FromTransaction: outTxn, ToTransaction: inTxn}, nil
}

func (s *Service) AccountHistory(ctx context.Context, accountID int64) ([]domain.Transaction, error) {
	logger.Info("account transaction history request", logger.Fields{"account_id": accountID})
	if accountID == 0 {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.ListAccountTransactions(ctx, accountID)
}

func (s *Service) ClientHistory(ctx context.Context, clientID int64) ([]domain.Transaction, error) {
	logger.Info("client transaction history request", logger.Fields{"client_id": clientID})
	if clientID == 0 {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.ListClientTransactions(ctx, clientID)
}

func (s *Service) mapRepoError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "no rows") {
		return domain.ErrAccountNotFound
	}
	return fmt.Errorf("transaction service: %w", err)
}
