package repository

import (
	"context"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
)

type LoyaltyCommandRepository interface {
	RedeemPoints(ctx context.Context, userId string, points int) error
	CollectPoints(ctx context.Context, userId string, points int, product *model.Product) error
}

type LoyaltyQueryRepository interface {
	GetPoints(ctx context.Context, userId string) (int, error)
	GetTransactions(ctx context.Context, userId string) (*[]model.Transaction, error)
}
