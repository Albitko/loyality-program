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
)
