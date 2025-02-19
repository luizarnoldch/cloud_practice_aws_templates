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

func Handler(ctx context.Context) (string, error) {
	client, err := lib.NewDynamoDBClient(ctx)
	if err != nil {
		return "", err
	}

	user_repository := src.NewUserRepositoryDynamoDB(
		ctx, client, USER_TABLE)

	newUserInput := src.User{
		Name: "John Doe",
		Age:  16,
	}

	newUser, err := user_repository.SaveUser(&newUserInput)
	if err != nil {
		return "", err
	}
	marshaled_user, err := json.Marshal(newUser)
	if err != nil {
		return "", err
	}

	return string(marshaled_user), nil
}

func main() {
	lambda.Start(Handler)
}
