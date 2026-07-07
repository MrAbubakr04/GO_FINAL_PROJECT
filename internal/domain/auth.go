package domain

import "time"

type Phone struct {
	ID       int64
	PhoneNum string
	ClientID *int64
	ActiveTo *time.Time
	Status   string
}

type Account struct {
	ID         int64
	PhoneNum   string
	PINHash    string
	BalanceTJ  int64
	BalanceRU  int64
	BalanceEN  int64
	Device     string
	IsActive   bool
	StatusID   int
	StatusCode string
	ActiveTo   *time.Time
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
