package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"os"
	// registering database driver
	_ "github.com/lib/pq"

	"github.com/rs/zerolog/log"
)

// New returns a new instance of Data with the database connection ready.
func NewDynamoDB() (*Service, error) {
	db, err := getConnectionDynamo()
	if err != nil {
		return nil, err
	}

	return &Service{DB: db}, nil
}

// Service is struct for library database/sql
type Service struct {
	DB *dynamodb.DynamoDB
}

func getConnectionDynamo() (*dynamodb.DynamoDB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	DbHost := fmt.Sprintf(os.Getenv("DB_HOST")) //"127.0.0.1"
	config := &aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Endpoint:    &DbHost,
		Credentials: credentials.NewEnvCredentials(),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)
	log.Info().Msg("Connected to database")

	// Verifica si la tabla ya existe
	_, err = svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String("loyalty"),
	})
	if err == nil {
		log.Info().Msg("Table exists! return instance")
		// La tabla ya existe, simplemente retorna una instancia del LoyaltyDB
		return svc, nil
	}

	err = CreateTable(svc)
	if err != nil {
		log.Info().Msg("Error on create table")
		log.Info().Msg(err.Error())
		return nil, err
	}
	log.Info().Msg("Created table")
	return svc, nil
}

func CreateTable(ser *dynamodb.DynamoDB) error {

	input := &dynamodb.CreateTableInput{
		TableName: aws.String("loyalty"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("UserId"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("transactionId"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("UserId"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("transactionId"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("points"), AttributeType: aws.String("S")},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),

		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("PointsIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("UserId"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("points"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(8),
					WriteCapacityUnits: aws.Int64(8),
				},
			},
		},
	}

	_, err := ser.CreateTable(input)
	if err != nil {
		return err
	}

	return nil
}
