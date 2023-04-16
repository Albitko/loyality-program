package entities

import (
	"errors"
)

var (
	ErrLoginAlreadyInUse = errors.New("login already exists")
)
