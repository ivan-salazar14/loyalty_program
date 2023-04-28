package repository

import (
	"context"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/model"
)

type SaleRepository interface {
	CreateSale(ctx context.Context, sale *model.Sale) (string, error)
	GetSale(ctx context.Context, SaleId string) (*model.Sale, error)
}
