package postgres

import (
	"context"
	"fmt"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetPhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error) {
	var phone domain.Phone
	err := r.db.QueryRow(ctx, `
		SELECT id, phone_num, client_id, active_to
		FROM phones
		WHERE phone_num = $1
	`, phoneNum).Scan(&phone.ID, &phone.PhoneNum, &phone.ClientID, &phone.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get phone by number: %w", err)
	}
	return &phone, nil
}

func (r *TransactionRepository) GetAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
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
		return nil, fmt.Errorf("get account by phone: %w", err)
	}
	return &account, nil
}

func (r *TransactionRepository) GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(ctx, `
		SELECT a.id, a.phone_num, a.pin, a.balance_tj, a.balance_ru, a.balance_en, a.device, a.is_active, a.status_id, s.code, a.active_to
		FROM accounts a
		LEFT JOIN account_statuses s ON s.id = a.status_id
		WHERE a.id = $1 AND a.active_to IS NULL
	`, accountID).Scan(&account.ID, &account.PhoneNum, &account.PINHash, &account.BalanceTJ, &account.BalanceRU, &account.BalanceEN, &account.Device, &account.IsActive, &account.StatusID, &account.StatusCode, &account.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account by id: %w", err)
	}
	return &account, nil
}

func (r *TransactionRepository) GetAccountForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*domain.Account, error) {
	var account domain.Account
	err := tx.QueryRow(ctx, `
		SELECT a.id, a.phone_num, a.pin, a.balance_tj, a.balance_ru, a.balance_en, a.device, a.is_active, a.status_id, s.code, a.active_to
		FROM accounts a
		LEFT JOIN account_statuses s ON s.id = a.status_id
		WHERE a.id = $1 AND a.active_to IS NULL
		FOR UPDATE OF a
	`, accountID).Scan(&account.ID, &account.PhoneNum, &account.PINHash, &account.BalanceTJ, &account.BalanceRU, &account.BalanceEN, &account.Device, &account.IsActive, &account.StatusID, &account.StatusCode, &account.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account for update: %w", err)
	}
	return &account, nil
}

func (r *TransactionRepository) GetClientByID(ctx context.Context, clientID int64) (*domain.Client, error) {
	var client domain.Client
	err := r.db.QueryRow(ctx, `
		SELECT id, name, surname, fathername, doc_num, tin, birth_date, gender, address, active_to, dt_created, dt_updated
		FROM clients
		WHERE id = $1 AND active_to IS NULL
	`, clientID).Scan(&client.ID, &client.Name, &client.Surname, &client.Fathername, &client.DocNum, &client.TIN, &client.BirthDate, &client.Gender, &client.Address, &client.ActiveTo, &client.DTCreated, &client.DTUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get client by id: %w", err)
	}
	return &client, nil
}

func (r *TransactionRepository) IsAccountIdentified(ctx context.Context, tx pgx.Tx, accountID int64) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM accounts a
			LEFT JOIN phones p ON p.phone_num = a.phone_num
			LEFT JOIN clients c ON c.id = p.client_id
			WHERE a.id = $1 AND a.active_to IS NULL AND p.active_to IS NULL AND c.active_to IS NULL AND p.client_id IS NOT NULL
		)
	`, accountID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is account identified: %w", err)
	}
	return exists, nil
}

func (r *TransactionRepository) UpdateAccountBalance(ctx context.Context, tx pgx.Tx, accountID int64, delta int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE accounts
		SET balance_tj = balance_tj + $1
		WHERE id = $2 AND active_to IS NULL
	`, delta, accountID)
	return err
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, txn domain.Transaction) (*domain.Transaction, error) {
	var created domain.Transaction
	err := tx.QueryRow(ctx, `
		INSERT INTO transactions (from_acc_id, to_acc_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, from_acc_id, to_acc_id, amount, dt_created
	`, txn.FromAccountID, txn.ToAccountID, txn.Amount).
		Scan(&created.ID, &created.FromAccountID, &created.ToAccountID, &created.Amount, &created.DTCreated)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	return &created, nil
}

func (r *TransactionRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *TransactionRepository) ListAccountTransactions(ctx context.Context, accountID int64) ([]domain.Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, from_acc_id, to_acc_id, amount, dt_created
		FROM transactions
		WHERE from_acc_id = $1 OR to_acc_id = $1
		ORDER BY dt_created DESC
	`, accountID)
	if err != nil {
		return nil, fmt.Errorf("list account transactions: %w", err)
	}
	defer rows.Close()
	transactions := make([]domain.Transaction, 0)
	for rows.Next() {
		var txn domain.Transaction
		if err := rows.Scan(&txn.ID, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.DTCreated); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func (r *TransactionRepository) ListClientTransactions(ctx context.Context, clientID int64) ([]domain.Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.from_acc_id, t.to_acc_id, t.amount, t.dt_created
		FROM transactions t
		JOIN accounts a ON a.id = t.from_acc_id OR a.id = t.to_acc_id
		JOIN phones p ON p.phone_num = a.phone_num
		WHERE p.client_id = $1
		ORDER BY t.dt_created DESC
	`, clientID)
	if err != nil {
		return nil, fmt.Errorf("list client transactions: %w", err)
	}
	defer rows.Close()
	transactions := make([]domain.Transaction, 0)
	for rows.Next() {
		var txn domain.Transaction
		if err := rows.Scan(&txn.ID, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.DTCreated); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}
