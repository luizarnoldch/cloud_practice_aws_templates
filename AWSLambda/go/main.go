package main

import (
	"main/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.SimpleHandler)
}
