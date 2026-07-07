package repository

import (
	"context"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	GetActivePhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error)
	GetActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error)
	GetStatusIDByCode(ctx context.Context, code string) (int, error)
	GetStatusCodeByID(ctx context.Context, id int) (string, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreatePhone(ctx context.Context, tx pgx.Tx, phone domain.Phone) (*domain.Phone, error)
	CreateAccount(ctx context.Context, tx pgx.Tx, account domain.Account) (*domain.Account, error)
}

type ClientRepository interface {
	GetClientByTIN(ctx context.Context, tin string) (*domain.Client, error)
	GetClientByID(ctx context.Context, id int64) (*domain.Client, error)
	CreateClient(ctx context.Context, tx pgx.Tx, client domain.Client) (*domain.Client, error)
	UpdateClient(ctx context.Context, tx pgx.Tx, client domain.Client) (*domain.Client, error)
	DeactivateClient(ctx context.Context, tx pgx.Tx, clientID int64) error
	GetActivePhoneByNumberAndClient(ctx context.Context, phoneNum string, clientID int64) (*domain.Phone, error)
	GetActivePhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error)
	GetActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error)
	GetActiveAccountsByClientID(ctx context.Context, clientID int64) ([]domain.Account, error)
	DeactivateAccount(ctx context.Context, tx pgx.Tx, accountID int64) error
	DeactivatePhone(ctx context.Context, tx pgx.Tx, phoneID int64) error
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CreateUser(ctx context.Context, tx pgx.Tx, user domain.User) (*domain.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type TransactionRepository interface {
	GetPhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error)
	GetAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error)
	GetAccountForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*domain.Account, error)
	GetClientByID(ctx context.Context, clientID int64) (*domain.Client, error)
	IsAccountIdentified(ctx context.Context, tx pgx.Tx, accountID int64) (bool, error)
	UpdateAccountBalance(ctx context.Context, tx pgx.Tx, accountID int64, delta int64) error
	CreateTransaction(ctx context.Context, tx pgx.Tx, txn domain.Transaction) (*domain.Transaction, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	ListAccountTransactions(ctx context.Context, accountID int64) ([]domain.Transaction, error)
	ListClientTransactions(ctx context.Context, clientID int64) ([]domain.Transaction, error)
}
