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
	return svc, nil
}
