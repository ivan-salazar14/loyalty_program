package service

import (
	"context"
	"fmt"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/model"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/repository"
)

type SaleService struct {
	SaleRepository repository.SaleRepository
}

func NewSaleService(repo repository.SaleRepository) *SaleService {
	return &SaleService{
		SaleRepository: repo,
	}
}
func (c SaleService) CreateSale(ctx context.Context, sale *model.Sale) (string, error) {

	res, err := c.SaleRepository.CreateSale(ctx, sale)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c SaleService) GetSale(ctx context.Context, SaleId string) (*model.Sale, error) {
	sale, err := c.SaleRepository.GetSale(ctx, SaleId)
	fmt.Sprint("err en handler", err)
	if err != nil {
		return nil, err
	}
	return sale, nil

}
