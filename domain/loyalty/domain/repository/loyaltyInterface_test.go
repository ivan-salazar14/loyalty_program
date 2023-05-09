package repository

import (
	"context"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/mocks"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoyaltyCommandRepository_RedeemPoints(t *testing.T) {
	// setup
	ctx := context.Background()
	repo := &mocks.MockCommandHandler{}
	userId := "user123"
	points := 10
	repo.On("RedeemPoints", ctx, userId, points).Return(nil)

	// execution
	err := repo.RedeemPoints(ctx, userId, points)

	// validation
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestLoyaltyCommandRepository_CollectPoints(t *testing.T) {
	// setup
	ctx := context.Background()
	repo := &mocks.MockCommandHandler{}
	userId := "user123"
	points := 10
	product := &model.Product{ProductName: "Product 1", Price: 100}
	repo.On("CollectPoints", ctx, userId, product).Return(nil)

	// execution
	err := repo.CollectPoints(ctx, userId, points, product)

	// validation
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestLoyaltyQueryRepository_GetPoints(t *testing.T) {
	// setup
	ctx := context.Background()
	repo := &mocks.MockQueryHandler{}
	userId := "user123"
	userPoints := 50
	repo.On("GetPoints", ctx, userId).Return(userPoints, nil)

	// execution
	points, err := repo.GetPoints(ctx, userId)

	// validation
	assert.Nil(t, err)
	assert.Equal(t, userPoints, points)
	repo.AssertExpectations(t)
}

func TestLoyaltyQueryRepository_GetTransactions(t *testing.T) {
	// setup
	ctx := context.Background()
	repo := &mocks.MockQueryHandler{}

	userId := "user123"
	transactions := &[]model.Transaction{
		{TransactionId: "TRANSACTION#1", UserID: "user123", Points: "10", Type: "COLLECT", Product: "Product 1"},
		{TransactionId: "TRANSACTION#2", UserID: "user123", Points: "5", Type: "REDEEM"},
	}
	repo.On("GetTransactions", ctx, userId).Return(transactions, nil)

	// execution
	result, err := repo.GetTransactions(ctx, userId)

	// validation
	assert.Nil(t, err)
	assert.Equal(t, transactions, result)
	repo.AssertExpectations(t)
}
