package domain

import "time"

type Transaction struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	DTCreated     time.Time
}

type TransferResult struct {
	FromTransaction *Transaction
	ToTransaction   *Transaction
}
