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
