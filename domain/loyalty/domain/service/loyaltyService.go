package service

import (
	"context"
	"fmt"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/repository"
	"math"
	"os"
	"strconv"
)

type LoyaltyService struct {
	LoyaltyRepository      repository.LoyaltyCommandRepository
	LoyaltyQueryRepository repository.LoyaltyQueryRepository
}

func NewLoyaltyService(repo repository.LoyaltyCommandRepository, query repository.LoyaltyQueryRepository) *LoyaltyService {
	return &LoyaltyService{
		LoyaltyRepository:      repo,
		LoyaltyQueryRepository: query,
	}
}
func (c LoyaltyService) RedeemPoints(ctx context.Context, userId string, points int) error {

	err := c.LoyaltyRepository.RedeemPoints(ctx, userId, points)

	return err
}

func (c LoyaltyService) GetPoints(ctx context.Context, userId string) (int, error) {
	points, err := c.LoyaltyQueryRepository.GetPoints(ctx, userId)
	fmt.Sprint("err en handler", err)
	if err != nil {
		return -1, err
	}
	return points, nil

}

func (c LoyaltyService) GetTransactions(ctx context.Context, userId string) (*[]model.Transaction, error) {
	transactions, err := c.LoyaltyQueryRepository.GetTransactions(ctx, userId)
	fmt.Sprint("err en handler", err)
	if err != nil {
		return nil, err
	}
	return transactions, nil

}
func (c LoyaltyService) CollectPoints(ctx context.Context, userId string, product *model.Product) error {

	pointsToAdd, err := strconv.ParseFloat(os.Getenv("PERCENT_POINTS"), 64)
	if err != nil {
		return err
	}
	pointsToAdd = pointsToAdd * product.Price

	err = c.LoyaltyRepository.CollectPoints(ctx, userId, int(math.Round(pointsToAdd)), product)
	return err
}
