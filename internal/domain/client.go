package domain

import "time"

type Client struct {
	ID         int64
	Name       string
	Surname    string
	Fathername string
	DocNum     string
	TIN        string
	BirthDate  time.Time
	Gender     string
	Address    string
	ActiveTo   *time.Time
	DTCreated  time.Time
	DTUpdated  *time.Time
}

type User struct {
	ID           int64
	Login        string
	PasswordHash string
	Role         string
	ActiveFrom   time.Time
	ActiveTo     *time.Time
	DTCreated    time.Time
	DTUpdated    *time.Time
}
