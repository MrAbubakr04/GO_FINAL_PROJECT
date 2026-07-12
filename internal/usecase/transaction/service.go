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
	FromPhone string `json:"from_phone"`
	ToPhone   string `json:"to_phone"`
	Amount    int64  `json:"amount"`
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
	input.FromPhone = normalizePhone(input.FromPhone)
	input.ToPhone = normalizePhone(input.ToPhone)
	logger.Info("transaction transfer started", logger.Fields{"from_phone": input.FromPhone, "to_phone": input.ToPhone, "amount": input.Amount})
	if input.Amount <= 0 {
		logger.Warning("invalid transfer amount", logger.Fields{"amount": input.Amount})
		return nil, domain.ErrInvalidAmount
	}
	if input.FromPhone == "" || input.ToPhone == "" {
		return nil, domain.ErrInvalidInput
	}
	if input.FromPhone == input.ToPhone {
		logger.Warning("same phone transfer rejected", logger.Fields{"phone": input.FromPhone})
		return nil, domain.ErrSameAccountTransfer
	}

	fromAcc, err := s.getActiveAccountByPhone(ctx, input.FromPhone)
	if err != nil {
		logger.Error("failed to get sender account for transfer", logger.Fields{"phone": input.FromPhone}, err)
		return nil, err
	}
	toAcc, err := s.getActiveAccountByPhone(ctx, input.ToPhone)
	if err != nil {
		logger.Error("failed to get receiver account for transfer", logger.Fields{"phone": input.ToPhone}, err)
		return nil, err
	}
	if fromAcc.ID == toAcc.ID {
		logger.Warning("same account transfer rejected", logger.Fields{"account_id": fromAcc.ID})
		return nil, domain.ErrSameAccountTransfer
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		logger.Error("failed to begin transfer transaction", logger.Fields{"from": fromAcc.ID, "to": toAcc.ID}, err)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	firstID, secondID := fromAcc.ID, toAcc.ID
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	firstAcc, err := s.repo.GetAccountForUpdate(ctx, tx, firstID)
	if err != nil {
		logger.Error("failed to lock first transfer account", logger.Fields{"account_id": firstID}, err)
		return nil, err
	}
	if firstAcc == nil {
		return nil, domain.ErrAccountNotFound
	}

	secondAcc, err := s.repo.GetAccountForUpdate(ctx, tx, secondID)
	if err != nil {
		logger.Error("failed to lock second transfer account", logger.Fields{"account_id": secondID}, err)
		return nil, err
	}
	if secondAcc == nil {
		return nil, domain.ErrAccountNotFound
	}

	if firstAcc.ID == fromAcc.ID {
		fromAcc = firstAcc
		toAcc = secondAcc
	} else {
		fromAcc = secondAcc
		toAcc = firstAcc
	}

	if fromAcc.ActiveTo != nil || !fromAcc.IsActive || fromAcc.StatusCode != "active" {
		return nil, domain.ErrAccountInactive
	}
	if toAcc.ActiveTo != nil || !toAcc.IsActive || toAcc.StatusCode != "active" {
		return nil, domain.ErrAccountInactive
	}
	if fromAcc.BalanceTJ < input.Amount {
		return nil, domain.ErrInsufficientFunds
	}
	isIdentified, err := s.repo.IsAccountIdentified(ctx, tx, fromAcc.ID)
	if err != nil {
		logger.Error("failed to verify sender identification", logger.Fields{"account_id": fromAcc.ID}, err)
		return nil, err
	}
	if !isIdentified {
		return nil, domain.ErrAccountNotIdentified
	}
	toIdentified, err := s.repo.IsAccountIdentified(ctx, tx, toAcc.ID)
	if err != nil {
		logger.Error("failed to verify receiver identification", logger.Fields{"account_id": toAcc.ID}, err)
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
		logger.Error("failed to commit transfer transaction", logger.Fields{"from": fromAcc.ID, "to": toAcc.ID}, err)
		return nil, err
	}
	logger.Info("transaction transfer completed", logger.Fields{"from": fromAcc.ID, "to": toAcc.ID, "amount": input.Amount})
	return &domain.TransferResult{FromTransaction: outTxn, ToTransaction: inTxn}, nil
}

func (s *Service) getActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
	phone, err := s.repo.GetPhoneByNumber(ctx, phoneNum)
	if err != nil {
		return nil, err
	}
	if phone == nil {
		return nil, domain.ErrPhoneNotFound
	}
	if phone.ActiveTo != nil {
		return nil, domain.ErrPhoneInactive
	}

	account, err := s.repo.GetAccountByPhone(ctx, phoneNum)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, domain.ErrAccountNotFound
	}
	if account.ActiveTo != nil || !account.IsActive || account.StatusCode != "active" {
		return nil, domain.ErrAccountInactive
	}
	return account, nil
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

func normalizePhone(phone string) string {
	return strings.TrimSpace(strings.ReplaceAll(phone, " ", ""))
}
