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

	checkOrderExistTests := []struct {
		name            string
		orderNumber     int
		userFromRequest string
		userFromDB      string
		errorFromDB     error
		expectedErr     error
	}{
		{
			name:            "CheckOrderExist: order already created by another user",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user2",
			errorFromDB:     nil,
			expectedErr:     entities.ErrOrderAlreadyCreatedByAnotherUser,
		},
		{
			name:            "CheckOrderExist: order already created by the same user",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user1",
			errorFromDB:     nil,
			expectedErr:     entities.ErrOrderAlreadyCreatedByThisUser,
		},
		{
			name:            "CheckOrderExist: internal error",
			orderNumber:     1234567,
			userFromRequest: "user1",
			userFromDB:      "user1",
			errorFromDB:     errors.New("internal error"),
			expectedErr:     errors.New("internal error"),
		},
	}
	for _, tt := range checkOrderExistTests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrdersRepository.EXPECT().
				GetUserForOrder(ctx, strconv.Itoa(tt.orderNumber)).
				Return(tt.userFromDB, tt.errorFromDB).
				Once()
			err := ordersProcessor.CheckOrderExist(ctx, tt.orderNumber, tt.userFromRequest)
			assert.Equal(t, tt.expectedErr, err)
		})
	}

	getUserOrderTests := []struct {
		name        string
		userID      string
		errorFromDB error
		expectedErr error
	}{
		{
			name:        "CheckOrderExist: return orders without errors",
			userID:      "1234567",
			errorFromDB: nil,
			expectedErr: nil,
		},
		{
			name:        "CheckOrderExist: return orders with errors",
			userID:      "1234567",
			errorFromDB: errors.New("error from DB"),
			expectedErr: errors.New("error from DB"),
		},
	}
	for _, tt := range getUserOrderTests {
		t.Run(tt.name, func(t *testing.T) {
			var userOrdersFromDB []entities.OrderWithTime
			mockOrdersRepository.EXPECT().
				GetOrdersForUser(ctx, tt.userID).
				Return(userOrdersFromDB, tt.errorFromDB).
				Once()
			_, err := ordersProcessor.GetUserOrder(ctx, tt.userID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}

	registerOrderTests := []struct {
		name        string
		orderNumber int
		userID      string
		errorFromDB error
		expectedErr error
	}{
		{
			name:        "RegisterOrder: return orders without errors",
			orderNumber: 111111,
			userID:      "1234567",
			errorFromDB: nil,
			expectedErr: nil,
		},
		{
			name:        "RegisterOrder: return orders with errors",
			orderNumber: 111111,
			userID:      "1234567",
			errorFromDB: errors.New("error from DB"),
			expectedErr: errors.New("error from DB"),
		},
	}
	for _, tt := range registerOrderTests {
		t.Run(tt.name, func(t *testing.T) {
			var order entities.Order
			order.OrderID = strconv.Itoa(tt.orderNumber)
			order.Status = "NEW"

			mockOrdersRepository.EXPECT().
				CreateOrder(ctx, order, tt.userID).
				Return(tt.errorFromDB).
				Once()
			mockOrdersQueue.EXPECT().
				Push(order)
			err := ordersProcessor.RegisterOrder(ctx, tt.orderNumber, tt.userID)
			assert.Equal(t, tt.expectedErr, err)

		})
	}
}
