package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Albitko/loyalty-program/internal/entities"
)

func TestUserAuthenticator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	mockUserRepository := newMockUserRepository(t)
	userAuthenticator := NewAuthenticator(mockUserRepository, "secret")

	createAccessTokenTests := []struct {
		name          string
		user          entities.User
		expectedToken string
		expectedErr   error
	}{
		{
			name: "CreateAccessToken: success",
			user: entities.User{
				ID:       "12345",
				Login:    "admin",
				Password: "pass",
			},
			expectedToken: "",
			expectedErr:   nil,
		},
	}
	for _, tt := range createAccessTokenTests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := userAuthenticator.CreateAccessToken(tt.user)
			assert.Equal(t, tt.expectedErr, err)
			assert.NotEqual(t, tt.expectedToken, token)
		})
	}

	AuthTests := []struct {
		name         string
		login        string
		password     string
		userFromDB   entities.User
		errFromDB    error
		expectedUser entities.User
		expectedErr  error
	}{
		{
			name:     "Auth: success",
			login:    "login",
			password: "<password>",
			userFromDB: entities.User{
				ID:       "123456",
				Login:    "login",
				Password: "dd81ca61fb57a4ff454c1cf89335a1f5e96afa849dfad4e0116b6ec35309fdea",
			},
			errFromDB: nil,
			expectedUser: entities.User{
				ID:       "123456",
				Login:    "login",
				Password: "dd81ca61fb57a4ff454c1cf89335a1f5e96afa849dfad4e0116b6ec35309fdea",
			},
			expectedErr: nil,
		},
		{
			name:     "Auth: invalid credentials",
			login:    "login",
			password: "<password>",
			userFromDB: entities.User{
				ID:       "123456",
				Login:    "login",
				Password: "qwertyuio12345678",
			},
			errFromDB: nil,
			expectedUser: entities.User{
				ID:       "123456",
				Login:    "login",
				Password: "qwertyuio12345678",
			},
			expectedErr: entities.ErrInvalidCredentials,
		},
		{
			name:         "Auth: DB error",
			login:        "login",
			password:     "<password>",
			userFromDB:   entities.User{},
			errFromDB:    errors.New("database error"),
			expectedUser: entities.User{},
			expectedErr:  errors.New("database error"),
		},
	}
	for _, tt := range AuthTests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository.EXPECT().
				GetCredentials(ctx, tt.login).
				Return(tt.userFromDB, tt.errFromDB).
				Once()
			user, err := userAuthenticator.Auth(ctx, tt.login, tt.password)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, user)
		})
	}

}
