package service

import (
	"context"
	"fmt"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/repository"
	"strconv"
)

type LoyaltyService struct {
	LoyaltyRepository repository.LoyaltyRepository
}

func NewLoyaltyService(repo repository.LoyaltyRepository) *LoyaltyService {
	return &LoyaltyService{
		LoyaltyRepository: repo,
	}
}
func (c LoyaltyService) RedeemPoints(ctx context.Context, userId string, points int) error {

	err := c.LoyaltyRepository.RedeemPoints(ctx, userId, points)

	return err
}

func (c LoyaltyService) GetPoints(ctx context.Context, userId string) (int, error) {
	points, err := c.LoyaltyRepository.GetPoints(ctx, userId)
	fmt.Sprint("err en handler", err)
	if err != nil {
		return -1, err
	}
	return points, nil

}
func (c LoyaltyService) CollectPoints(ctx context.Context, userId string, points string) error {

	pointsToCollet, err := strconv.Atoi(points)
	if err != nil {
		fmt.Println("Error during conversion")
		return err
	}
	err = c.LoyaltyRepository.CollectPoints(ctx, userId, pointsToCollet)
	return err
}
