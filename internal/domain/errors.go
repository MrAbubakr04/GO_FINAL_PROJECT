package domain

import "errors"

var (
	ErrPhoneNotFound        = errors.New("phone not found")
	ErrPhoneAlreadyExists   = errors.New("phone already exists")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrPhoneInactive        = errors.New("phone inactive")
	ErrInvalidPin           = errors.New("invalid pin")
	ErrInvalidInput         = errors.New("invalid input")
)
