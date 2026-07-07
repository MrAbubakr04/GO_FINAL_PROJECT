package postgres

import (
	"context"
	"fmt"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetActivePhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error) {
	var phone domain.Phone
	err := r.db.QueryRow(ctx, `
		SELECT id, phone_num, client_id, dt_created, dt_updated, active_to
		FROM phones
		WHERE phone_num = $1 AND active_to IS NULL
	`, phoneNum).Scan(&phone.ID, &phone.PhoneNum, &phone.ClientID, new(interface{}), new(interface{}), &phone.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get active phone: %w", err)
	}
	return &phone, nil
}

func (r *AuthRepository) GetActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(ctx, `
		SELECT a.id, a.phone_num, a.pin, a.balance_tj, a.balance_ru, a.balance_en, a.device, a.is_active, a.status_id, s.code, a.active_to
		FROM accounts a
		LEFT JOIN account_statuses s ON s.id = a.status_id
		WHERE a.phone_num = $1 AND a.active_to IS NULL
	`, phoneNum).Scan(&account.ID, &account.PhoneNum, &account.PINHash, &account.BalanceTJ, &account.BalanceRU, &account.BalanceEN, &account.Device, &account.IsActive, &account.StatusID, &account.StatusCode, &account.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get active account: %w", err)
	}
	return &account, nil
}

func (r *AuthRepository) GetStatusIDByCode(ctx context.Context, code string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM account_statuses WHERE code = $1`, code).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("get status id: %w", err)
	}
	return id, nil
}

func (r *AuthRepository) GetStatusCodeByID(ctx context.Context, id int) (string, error) {
	var code string
	err := r.db.QueryRow(ctx, `SELECT code FROM account_statuses WHERE id = $1`, id).Scan(&code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get status code: %w", err)
	}
	return code, nil
}

func (r *AuthRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *AuthRepository) CreatePhone(ctx context.Context, tx pgx.Tx, phone domain.Phone) (*domain.Phone, error) {
	var created domain.Phone
	err := tx.QueryRow(ctx, `
		INSERT INTO phones (phone_num, client_id) VALUES ($1, $2) RETURNING id, phone_num, client_id, active_to
	`, phone.PhoneNum, phone.ClientID).Scan(&created.ID, &created.PhoneNum, &created.ClientID, &created.ActiveTo)
	if err != nil {
		return nil, fmt.Errorf("create phone: %w", err)
	}
	return &created, nil
}

func (r *AuthRepository) CreateAccount(ctx context.Context, tx pgx.Tx, account domain.Account) (*domain.Account, error) {
	var created domain.Account
	err := tx.QueryRow(ctx, `
		INSERT INTO accounts (phone_num, pin, balance_tj, balance_ru, balance_en, device, is_active, status_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, phone_num, pin, balance_tj, balance_ru, balance_en, device, is_active, status_id, active_to
	`, account.PhoneNum, account.PINHash, account.BalanceTJ, account.BalanceRU, account.BalanceEN, account.Device, account.IsActive, account.StatusID).
		Scan(&created.ID, &created.PhoneNum, &created.PINHash, &created.BalanceTJ, &created.BalanceRU, &created.BalanceEN, &created.Device, &created.IsActive, &created.StatusID, &created.ActiveTo)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return &created, nil
}
