package usecase

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Albitko/loyalty-program/internal/entities"
)

func TestOrdersProcessor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	mockOrdersRepository := newMockOrdersRepository(t)
	mockOrdersQueue := newMockOrdersQueue(t)

	ordersProcessor := NewOrdersProcessor(mockOrdersRepository, mockOrdersQueue)

	// Test `CheckOrderExist`
	tests := []struct {
		name            string
		orderNumber     int
		userFromRequest string
		userFromDB      string
		errorFromDB     error
		expectedErr     error
	}{
		{
			name:            "order already created by another user",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user2",
			errorFromDB:     nil,
			expectedErr:     entities.ErrOrderAlreadyCreatedByAnotherUser,
		},
		{
			name:            "order already created by another user",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user1",
			errorFromDB:     nil,
			expectedErr:     entities.ErrOrderAlreadyCreatedByThisUser,
		},
		{
			name:            "internal error",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user1",
			errorFromDB:     errors.New("internal error"),
			expectedErr:     errors.New("internal error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrdersRepository.EXPECT().
				GetUserForOrder(ctx, strconv.Itoa(tt.orderNumber)).
				Return(tt.userFromDB, tt.errorFromDB).
				Once()
			err := ordersProcessor.CheckOrderExist(ctx, tt.orderNumber, tt.userFromRequest)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
