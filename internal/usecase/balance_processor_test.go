package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Albitko/loyalty-program/internal/entities"
)

func TestBalanceProcessor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	mockBalanceRepository := newMockBalanceRepository(t)

	balanceProcessor := NewBalanceProcessor(mockBalanceRepository)

	getUserBalanceTests := []struct {
		name             string
		userID           string
		accrualsTotal    float64
		withdrawalsTotal float64
		errFromDB        error
		expectedBalance  entities.Balance
		expectedErr      error
	}{
		{
			name:             "GetUserBalance: positive",
			userID:           "123456",
			accrualsTotal:    1000.6,
			withdrawalsTotal: 345,
			errFromDB:        nil,
			expectedBalance: entities.Balance{
				Current:   655.6,
				Withdrawn: 345,
			},
			expectedErr: nil,
		},
	}
	for _, tt := range getUserBalanceTests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepository.EXPECT().
				GetUserBalance(ctx, tt.userID).
				Return(tt.accrualsTotal, tt.errFromDB).
				Once()
			mockBalanceRepository.EXPECT().
				GetUserWithdrawn(ctx, tt.userID).
				Return(tt.withdrawalsTotal, tt.errFromDB).
				Once()
			balance, err := balanceProcessor.GetUserBalance(ctx, tt.userID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedBalance, balance)
		})
	}

	withdrawTests := []struct {
		name            string
		userID          string
		withdrawRequest entities.Withdraw
		balanceFromDB   entities.Balance
		errFromDB       error
		expectedErr     error
	}{
		{
			name:   "Withdraw: positive",
			userID: "123456",
			withdrawRequest: entities.Withdraw{
				Order: "777777777",
				Sum:   1000,
			},
			balanceFromDB: entities.Balance{
				Current:   1000,
				Withdrawn: 0,
			},
			errFromDB:   nil,
			expectedErr: nil,
		},
	}
	for _, tt := range withdrawTests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepository.EXPECT().
				GetUserBalance(ctx, tt.userID).
				Return(tt.balanceFromDB.Current, tt.errFromDB).
				Once()
			mockBalanceRepository.EXPECT().
				GetUserWithdrawn(ctx, tt.userID).
				Return(tt.balanceFromDB.Withdrawn, tt.errFromDB).
				Once()
			mockBalanceRepository.EXPECT().
				Withdraw(ctx, tt.userID, tt.withdrawRequest).
				Return(tt.errFromDB).
				Once()
			err := balanceProcessor.Withdraw(ctx, tt.userID, tt.withdrawRequest)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
