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
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("userId %s " + userId)
	userPoints, err := db.GetPoints(ctx, userId)

	if err != nil {
		return err
	}

	// Check if the user has enough points to redeem
	if userPoints < points {
		return fmt.Errorf("user does not have enough points to redeem %d", userPoints)
	}

	// Create a new PutItemInput object to insert a new item into the table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(fmt.Sprintf("USER#%s", userId)),
			},
			"transactionId": {
				S: aws.String(fmt.Sprintf("TRANSACTION#%d", now)),
			},
			"points": {
				S: aws.String(strconv.Itoa(-points)),
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
		log.Log().Msg(err.Error())
		return err
	}
	err = db.updatePoints(ctx, userId, -points)

	if err != nil {
		log.Log().Msg(err.Error())
		return err
	}
	return nil

}

func (db *dbLoyalty) GetPoints(ctx context.Context, userId string) (int, error) {

	var points struct {
		UserId        string
		Points        int    `DynamoDB:"points"`
		TransactionID string `DynamoDB:"transactionId"`
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId":        {S: aws.String("USER#" + userId)},
			"transactionId": {S: aws.String("TotalPoints")},
		},
	}
	result, err := db.ConnectDb.DB.GetItem(input)
	if err != nil {
		return 0, err
	}

	if result.Item == nil {
		fmt.Printf("total not found,lets create ")
		return 0, err
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &points); err != nil {
		return 0, err
	}
	fmt.Printf("total found --------- %s", result.Item)
	fmt.Printf("total  --------- %d", points.Points)
	fmt.Printf("UserId  --------- %s", points.UserId)
	log.Info().Msg(result.String())
	return points.Points, nil
}

func (db *dbLoyalty) GetTransactions(ctx context.Context, userId string) (*[]model.Transaction, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String("loyalty"),
		KeyConditionExpression: aws.String("UserId = :uid AND begins_with(transactionId, :tid)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":uid": {S: aws.String("USER#" + userId)},
			":tid": {S: aws.String("TRANSACTION#")},
		},
	}
	result, err := db.ConnectDb.DB.Query(input)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, err
	}
	// Unmarshal the results
	var transactions []model.Transaction
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &transactions)
	if err != nil {
		fmt.Printf("vacio UnmarshalListOfMaps --------- %s", err)
		return nil, err
	}

	return &transactions, nil
}

func (db *dbLoyalty) CollectPoints(ctx context.Context, userId string, points int, product *model.Product) error {

	now := time.Now().Unix()
	// Create a new Transaction object. insert a new item into the table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(fmt.Sprintf("USER#%s", userId)),
			},
			"transactionId": {
				S: aws.String(fmt.Sprintf("TRANSACTION#%d", now)),
			},

			"points": {
				S: aws.String(strconv.Itoa(points)),
			},
			"type_tran": {
				S: aws.String("COLLECT"),
			},
			"product": {
				S: aws.String(product.ProductName),
			},
			/*	"price": {
				N: aws.String(product.Price),
			},*/
			"date": {
				S: aws.String(time.Unix(now, 0).Format(time.RFC3339)),
			},
		},
	}
	_, err := db.ConnectDb.DB.PutItemWithContext(ctx, input)
	if err != nil {
		log.Info().Msg("entroerror put transact err.Error()")
		return err
	}

	err = db.updatePoints(ctx, userId, points)
	if err != nil {
		log.Info().Msg("entroerror updatepoints err.Error()")

		return err
	}
	return nil
}

func (db *dbLoyalty) updatePoints(ctx context.Context, userId string, points int) error {
	key := map[string]*dynamodb.AttributeValue{
		"UserId": {
			S: aws.String("USER#" + userId),
		},
		"transactionId": {
			S: aws.String("TotalPoints"),
		},
	}
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String("loyalty"),
		Key:              key,
		UpdateExpression: aws.String("SET points = if_not_exists(points, :zero) + :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":zero": {
				N: aws.String("0"),
			},
			":val": {
				N: aws.String(strconv.FormatInt(int64(points), 10)),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}
	_, err := db.ConnectDb.DB.UpdateItem(input)
	if err != nil {
		return err
	}
	fmt.Printf("updated balance  ")
	return nil
}

func NewConnection(connectDb *database.Service) repoDomain.LoyaltyCommandRepository {
	return &dbLoyalty{
		ConnectDb: connectDb,
		tableName: os.Getenv("TABLE_NAME"),
	}
}
func NewConnectionQuery(connectDb *database.Service) repoDomain.LoyaltyQueryRepository {
	return &dbLoyalty{
		ConnectDb: connectDb,
		tableName: os.Getenv("TABLE_NAME"),
	}
}
