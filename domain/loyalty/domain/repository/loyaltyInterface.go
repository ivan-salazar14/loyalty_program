package repository

import (
	"context"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
)

type LoyaltyRepository interface {
	RedeemPoints(ctx context.Context, userId string, points int) error
	GetPoints(ctx context.Context, userId string) (int, error)
	CollectPoints(ctx context.Context, userId string, points int, product *model.Product) error
	GetTransactions(ctx context.Context, userId string) (*[]model.Transaction, error)
}
