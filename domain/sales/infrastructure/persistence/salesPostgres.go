package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/model"
	repoDomain "github.com/ivan-salazar14/firstGoPackage/domain/sales/domain/repository"
	"github.com/ivan-salazar14/firstGoPackage/infrastructure/database"
)

const (
	Select = "SELECT * FROM public.sales WHERE sale_id = $1"
	Insert = "INSERT INTO public.sales  (product_id, user_id, price) VALUES ($1,$2,$3)"
)

type dbSale struct {
	ConnectDb *database.DataDB
}

func NewConnection(connectDb *database.DataDB) repoDomain.SaleRepository {
	return &dbSale{
		ConnectDb: connectDb,
	}
}
func (db *dbSale) CreateSale(ctx context.Context, sale *model.Sale) (string, error) {
	var idResult string
	smt, err := db.ConnectDb.DB.PrepareContext(ctx, Insert)

	if err != nil {
		return "", err
	}

	defer smt.Close()

	row := smt.QueryRowContext(ctx, &sale.ProductId, &sale.CustomerId, &sale.Price)

	if err = row.Scan(&idResult); err != sql.ErrNoRows {
		return "", err
	}

	return idResult, nil
}

func (db *dbSale) GetSale(ctx context.Context, SaleId string) (*model.Sale, error) {
	fmt.Sprint("entro GetSale")
	smt, err := db.ConnectDb.DB.PrepareContext(ctx, Select)
	fmt.Sprint(smt)
	if err != nil {

		fmt.Sprint(err)
		return nil, err
	}
	defer smt.Close()
	row := smt.QueryRowContext(ctx, &SaleId)
	sale := model.Sale{}

	err = row.Scan(&sale.Id, &sale.Price, &sale.CustomerId, &sale.ProductId)

	if err != nil {
		fmt.Sprint(err)
		return nil, err
	}
	//TODO implement me
	return &sale, nil
}
