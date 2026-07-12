package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetClientByTIN(ctx context.Context, tin string) (*domain.Client, error) {
	var client domain.Client
	err := r.db.QueryRow(ctx, `
		SELECT id, name, surname, fathername, doc_num, tin, birth_date, gender, address, active_to, dt_created, dt_updated
		FROM clients
		WHERE tin = $1 AND active_to IS NULL
	`, tin).Scan(&client.ID, &client.Name, &client.Surname, &client.Fathername, &client.DocNum, &client.TIN, &client.BirthDate, &client.Gender, &client.Address, &client.ActiveTo, &client.DTCreated, &client.DTUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get client by tin: %w", err)
	}
	return &client, nil
}

func (r *ClientRepository) GetClientByID(ctx context.Context, id int64) (*domain.Client, error) {
	var client domain.Client
	err := r.db.QueryRow(ctx, `
		SELECT id, name, surname, fathername, doc_num, tin, birth_date, gender, address, active_to, dt_created, dt_updated
		FROM clients
		WHERE id = $1 AND active_to IS NULL
	`, id).Scan(&client.ID, &client.Name, &client.Surname, &client.Fathername, &client.DocNum, &client.TIN, &client.BirthDate, &client.Gender, &client.Address, &client.ActiveTo, &client.DTCreated, &client.DTUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get client by id: %w", err)
	}
	return &client, nil
}

func (r *ClientRepository) CreateClient(ctx context.Context, tx pgx.Tx, client domain.Client) (*domain.Client, error) {
	var created domain.Client
	err := tx.QueryRow(ctx, `
		INSERT INTO clients (name, surname, fathername, doc_num, tin, birth_date, gender, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, surname, fathername, doc_num, tin, birth_date, gender, address, active_to, dt_created, dt_updated
	`, client.Name, client.Surname, client.Fathername, client.DocNum, client.TIN, client.BirthDate, client.Gender, client.Address).
		Scan(&created.ID, &created.Name, &created.Surname, &created.Fathername, &created.DocNum, &created.TIN, &created.BirthDate, &created.Gender, &created.Address, &created.ActiveTo, &created.DTCreated, &created.DTUpdated)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return &created, nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, tx pgx.Tx, client domain.Client) (*domain.Client, error) {
	var updated domain.Client
	err := tx.QueryRow(ctx, `
		UPDATE clients
		SET name = $1, surname = $2, fathername = $3, doc_num = $4, tin = $5, birth_date = $6, gender = $7, address = $8, dt_updated = NOW()
		WHERE id = $9 AND active_to IS NULL
		RETURNING id, name, surname, fathername, doc_num, tin, birth_date, gender, address, active_to, dt_created, dt_updated
	`, client.Name, client.Surname, client.Fathername, client.DocNum, client.TIN, client.BirthDate, client.Gender, client.Address, client.ID).
		Scan(&updated.ID, &updated.Name, &updated.Surname, &updated.Fathername, &updated.DocNum, &updated.TIN, &updated.BirthDate, &updated.Gender, &updated.Address, &updated.ActiveTo, &updated.DTCreated, &updated.DTUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("update client: %w", err)
	}
	return &updated, nil
}

func (r *ClientRepository) DeactivateClient(ctx context.Context, tx pgx.Tx, clientID int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE clients
		SET active_to = NOW(), dt_updated = NOW()
		WHERE id = $1 AND active_to IS NULL
	`, clientID)
	return err
}

func (r *ClientRepository) GetActivePhoneByNumberAndClient(ctx context.Context, phoneNum string, clientID int64) (*domain.Phone, error) {
	var phone domain.Phone
	err := r.db.QueryRow(ctx, `
		SELECT id, phone_num, client_id, active_to
		FROM phones
		WHERE phone_num = $1 AND client_id = $2 AND active_to IS NULL
	`, phoneNum, clientID).Scan(&phone.ID, &phone.PhoneNum, &phone.ClientID, &phone.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get phone by number and client: %w", err)
	}
	return &phone, nil
}

func (r *ClientRepository) GetActivePhoneByNumber(ctx context.Context, phoneNum string) (*domain.Phone, error) {
	var phone domain.Phone
	err := r.db.QueryRow(ctx, `
		SELECT id, phone_num, client_id, active_to
		FROM phones
		WHERE phone_num = $1 AND active_to IS NULL
	`, phoneNum).Scan(&phone.ID, &phone.PhoneNum, &phone.ClientID, &phone.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get phone by number: %w", err)
	}
	return &phone, nil
}

func (r *ClientRepository) AttachPhoneToClient(ctx context.Context, tx pgx.Tx, phoneID int64, clientID int64) error {
	tag, err := tx.Exec(ctx, `
		UPDATE phones
		SET client_id = $1, dt_updated = NOW()
		WHERE id = $2 AND client_id IS NULL AND active_to IS NULL
	`, clientID, phoneID)
	if err != nil {
		return fmt.Errorf("attach phone to client: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPhoneAlreadyUsed
	}
	return nil
}

func (r *ClientRepository) GetActiveAccountByPhone(ctx context.Context, phoneNum string) (*domain.Account, error) {
	var account domain.Account
	err := r.db.QueryRow(ctx, `
		SELECT a.id, a.phone_num, a.pin, a.balance_tj, a.balance_ru, a.balance_en, a.device, a.is_active, a.status_id, a.active_to
		FROM accounts a
		LEFT JOIN phones p ON p.phone_num = a.phone_num
		WHERE a.phone_num = $1 AND a.active_to IS NULL AND p.active_to IS NULL
	`, phoneNum).Scan(&account.ID, &account.PhoneNum, &account.PINHash, &account.BalanceTJ, &account.BalanceRU, &account.BalanceEN, &account.Device, &account.IsActive, &account.StatusID, &account.ActiveTo)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account by phone: %w", err)
	}
	return &account, nil
}

func (r *ClientRepository) GetActiveAccountsByClientID(ctx context.Context, clientID int64) ([]domain.Account, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, phone_num, pin, balance_tj, balance_ru, balance_en, device, is_active, status_id, active_to
		FROM accounts
		WHERE phone_num IN (SELECT phone_num FROM phones WHERE client_id = $1 AND active_to IS NULL) AND active_to IS NULL
	`, clientID)
	if err != nil {
		return nil, fmt.Errorf("get accounts by client: %w", err)
	}
	defer rows.Close()
	accounts := make([]domain.Account, 0)
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(&acc.ID, &acc.PhoneNum, &acc.PINHash, &acc.BalanceTJ, &acc.BalanceRU, &acc.BalanceEN, &acc.Device, &acc.IsActive, &acc.StatusID, &acc.ActiveTo); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *ClientRepository) DeactivateAccount(ctx context.Context, tx pgx.Tx, accountID int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE accounts
		SET active_to = NOW(), dt_updated = NOW()
		WHERE id = $1 AND active_to IS NULL
	`, accountID)
	return err
}

func (r *ClientRepository) DeactivatePhone(ctx context.Context, tx pgx.Tx, phoneID int64) error {
	_, err := tx.Exec(ctx, `
		UPDATE phones
		SET active_to = NOW(), dt_updated = NOW()
		WHERE id = $1 AND active_to IS NULL
	`, phoneID)
	return err
}

func (r *ClientRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, role, active_from, active_to, dt_created, dt_updated
		FROM users
		WHERE login = $1 AND active_to IS NULL
	`, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Role, &user.ActiveFrom, &user.ActiveTo, &user.DTCreated, &user.DTUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return &user, nil
}

func (r *ClientRepository) CreateUser(ctx context.Context, tx pgx.Tx, user domain.User) (*domain.User, error) {
	var created domain.User
	err := tx.QueryRow(ctx, `
		INSERT INTO users (login, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, login, password_hash, role, active_from, active_to, dt_created, dt_updated
	`, user.Login, user.PasswordHash, user.Role).
		Scan(&created.ID, &created.Login, &created.PasswordHash, &created.Role, &created.ActiveFrom, &created.ActiveTo, &created.DTCreated, &created.DTUpdated)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &created, nil
}

func (r *ClientRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *ClientRepository) Now() time.Time { return time.Now().UTC() }
