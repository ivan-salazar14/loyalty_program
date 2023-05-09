package service

import (
	"context"
	_ "errors"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/mocks"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestController_RedeemPoints(t *testing.T) {
	ctx := context.Background()
	userId := "123"
	points := 10

	mockCommand := new(mocks.MockCommandHandler)
	mockCommand.On("RedeemPoints", ctx, userId, points).Return(nil)

	controller := &LoyaltyService{
		LoyaltyRepository: mockCommand,
	}

	err := controller.RedeemPoints(ctx, userId, points)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mockCommand.AssertCalled(t, "RedeemPoints", ctx, userId, points)
}

func TestController_CollectPoints(t *testing.T) {
	ctx := context.Background()
	userId := "123"
	product := &model.Product{
		ProductId:   "456",
		ProductName: "Test Product",
		Price:       150,
	}
	os.Setenv("PERCENT_POINTS", "0.05")
	defer os.Unsetenv("PERCENT_POINTS")

	mockCommand := new(mocks.MockCommandHandler)
	mockCommand.On("CollectPoints", ctx, userId, product).Return(nil)

	controller := &LoyaltyService{
		LoyaltyRepository:      mockCommand,
		LoyaltyQueryRepository: nil,
	}

	err := controller.CollectPoints(ctx, userId, product)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mockCommand.AssertCalled(t, "CollectPoints", ctx, userId, product)
}

func TestController_GetPoints(t *testing.T) {
	ctx := context.Background()
	userId := "123"
	expectedPoints := 50

	mockQuery := new(mocks.MockQueryHandler)
	mockQuery.On("GetPoints", ctx, userId).Return(expectedPoints, nil)

	controller := &LoyaltyService{
		LoyaltyQueryRepository: mockQuery,
	}

	points, err := controller.GetPoints(ctx, userId)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(points, expectedPoints) {
		t.Errorf("expected %v, but got %v", expectedPoints, points)
	}

	mockQuery.AssertCalled(t, "GetPoints", ctx, userId)
}

func TestController_GetTransactions(t *testing.T) {
	ctx := context.Background()
	userId := "test-user"

	// Creamos un mock del QueryHandler para simular la respuesta
	//mockQueryHandler := &mocks.QueryHandler{}
	mockQueryHandler := new(mocks.MockQueryHandler)
	date := time.Now()
	transactions := []model.Transaction{
		{TransactionId: "TRANSACTION#1", UserID: "test-user", Points: "100", Type: "REDEEM", Date: date},
		{TransactionId: "TRANSACTION#2", UserID: "test-user", Points: "50", Type: "COLLECT", Date: date},
	}
	mockQueryHandler.On("GetTransactions", ctx, userId).Return(&transactions, nil)

	// Creamos el Controller con el QueryHandler mockeado
	controller := &LoyaltyService{
		LoyaltyRepository:      nil,
		LoyaltyQueryRepository: mockQueryHandler,
	}

	// Ejecutamos la función GetTransactions
	resp, err := controller.GetTransactions(ctx, userId)
	if err != nil {
		t.Errorf("Error obteniendo transacciones: %v", err)
		return
	}

	// Verificamos que la respuesta obtenida sea la esperada
	expectedTransactions := []model.Transaction{
		{TransactionId: "TRANSACTION#1", UserID: "test-user", Points: "100", Type: "REDEEM", Date: date},
		{TransactionId: "TRANSACTION#2", UserID: "test-user", Points: "50", Type: "COLLECT", Date: date},
	}
	if !reflect.DeepEqual(*resp, expectedTransactions) {
		t.Errorf("Respuesta inesperada. Esperada: %v, obtenida: %v", expectedTransactions, *resp)
		return
	}

	// Verificamos que el QueryHandler haya sido llamado con los argumentos correctos
	mockQueryHandler.AssertCalled(t, "GetTransactions", ctx, userId)
}
