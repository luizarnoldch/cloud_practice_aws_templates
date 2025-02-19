package main

import (
	"context"
	"encoding/json"
	"main/crud_dynamodb/src"
	"main/lib"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	USER_TABLE = os.Getenv("USER_TABLE")
)

type Event struct {
	Id string `json:"id"`
}

func Handler(ctx context.Context, event Event) (string, error) {
	client, err := lib.NewDynamoDBClient(ctx)
	if err != nil {
		return "", err
	}

	user_repository := src.NewUserRepositoryDynamoDB(
		ctx, client, USER_TABLE)

	user, err := user_repository.GetUserById(event.Id)
	if err != nil {
		return "", err
	}

	marshaled_user, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	return string(marshaled_user), nil
}

func main() {
	lambda.Start(Handler)
}
