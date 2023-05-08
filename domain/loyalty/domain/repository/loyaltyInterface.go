package repository

import (
	"context"
)

type LoyaltyRepository interface {
	RedeemPoints(ctx context.Context, userId string, points int) error
	GetPoints(ctx context.Context, userId string) (int, error)
	CollectPoints(ctx context.Context, userId string, points int) error
}
