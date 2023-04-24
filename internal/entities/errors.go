package entities

import (
	"errors"
)

var (
	ErrLoginAlreadyInUse                = errors.New("login already exists")
	ErrInvalidCredentials               = errors.New("invalid credentials")
	ErrOrderAlreadyCreatedByThisUser    = errors.New("user has already created this order")
	ErrOrderAlreadyCreatedByAnotherUser = errors.New("user has already created another order")
	ErrNoOrderForUser                   = errors.New("there is no order for this user")
	ErrInsufficientFunds                = errors.New("insufficient funds for this user")
	ErrNoWithdrawals                    = errors.New("no withdrawals for this user")
)
