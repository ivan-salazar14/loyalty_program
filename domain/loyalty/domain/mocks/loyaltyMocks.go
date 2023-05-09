package mocks

import (
	"context"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockCommandHandler struct {
	mock.Mock
}

func (m *MockCommandHandler) RedeemPoints(ctx context.Context, userId string, points int) error {
	args := m.Called(ctx, userId, points)
	return args.Error(0)
}

func (m *MockCommandHandler) CollectPoints(ctx context.Context, userId string, points int, product *model.Product) error {
	args := m.Called(ctx, userId, product)
	return args.Error(0)
}

type MockQueryHandler struct {
	mock.Mock
}

func (m *MockQueryHandler) GetPoints(ctx context.Context, userId string) (int, error) {
	args := m.Called(ctx, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockQueryHandler) GetTransactions(ctx context.Context, userId string) (*[]model.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*[]model.Transaction), args.Error(1)
}
