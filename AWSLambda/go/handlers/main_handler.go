package handlers

import (
	"context"
	"log"
)

func SimpleHandler(ctx context.Context) (string, error) {
	log.Print("Lambda function executed successfully.")
	return "Hello from Lambda!", nil
}
