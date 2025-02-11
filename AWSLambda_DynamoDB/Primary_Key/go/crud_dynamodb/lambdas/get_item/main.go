package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func SimpleHandler(ctx context.Context) (string, error) {
	log.Print("Lambda function executed successfully.")
	return "Hello from Lambda!", nil
}

func main() {
	lambda.Start(SimpleHandler)
}

