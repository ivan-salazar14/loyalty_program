package persistence

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/model"
	repoDomain "github.com/ivan-salazar14/firstGoPackage/domain/loyalty/domain/repository"
	"github.com/ivan-salazar14/firstGoPackage/infrastructure/database"
	"os"
	"strconv"
	"time"
)

type dbLoyalty struct {
	tableName string
	ConnectDb *database.Service
}

func (db *dbLoyalty) RedeemPoints(ctx context.Context, userId string, points int) error {

	now := time.Now().Unix()

	userPoints, err := db.GetPoints(ctx, userId)
	if err != nil {
		return err
	}

	// Check if the user has enough points to redeem
	if userPoints < points {
		return fmt.Errorf("user does not have enough points to redeem")
	}

	// Create a new PutItemInput object to insert a new item into the table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(fmt.Sprintf("USER#%s", userId)),
			},
			"sk": {
				S: aws.String(fmt.Sprintf("POINT#%d", now)),
			},
			"userId": {
				S: aws.String(userId),
			},
			"points": {
				N: aws.String(strconv.Itoa(-points)),
			},
			"type_tran": {
				S: aws.String("REDEEM"),
			},
			"date": {
				S: aws.String(time.Unix(now, 0).Format(time.RFC3339)),
			},
		},
	}

	// Execute the PutItem request
	_, err = db.ConnectDb.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil

}

func (db *dbLoyalty) GetPoints(ctx context.Context, userId string) (int, error) {

	fmt.Sprint("entro GetPoints")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(userId)},
		},
	}
	result, err := db.ConnectDb.DB.GetItem(input)
	if err != nil {
		return 0, err
	}

	if result.Item == nil {
		return 0, fmt.Errorf("user no encontrado")
	}

	var user model.User
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return 0, err
	}
	return user.Points, nil
}

func (db *dbLoyalty) CollectPoints(ctx context.Context, userId string, points int) error {

	now := time.Now().Unix()

	// Create a new PutItemInput object to insert a new item into the table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(fmt.Sprintf("USER#%s", userId)),
			},
			"sk": {
				S: aws.String(fmt.Sprintf("POINT#%d", now)),
			},
			"userId": {
				S: aws.String(userId),
			},
			"points": {
				N: aws.String(strconv.Itoa(points)),
			},
			"type_tran": {
				S: aws.String("COLLECT"),
			},
			"date": {
				S: aws.String(time.Unix(now, 0).Format(time.RFC3339)),
			},
		},
	}

	// Execute the PutItem request
	_, err := db.ConnectDb.DB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func NewConnection(connectDb *database.Service) repoDomain.LoyaltyRepository {
	return &dbLoyalty{
		ConnectDb: connectDb,
		tableName: os.Getenv("TABLE_NAME"),
	}
}
