package lib

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DescribeLocalDynamoDBTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Printf("Table %s does not exists.", tableName)
			return notFoundErr
		}
		log.Printf("Unexpected error occurred while describing table %s: %s", tableName, err)
		return err
	}

	log.Printf("Table %s exists.", tableName)
	return nil
}

func DeleteLocalDynamoDBTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	err := DescribeLocalDynamoDBTable(ctx, client, tableName)

	if err != nil {
		return err
	}

	client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	log.Printf("Table %s deleted successfully", tableName)
	return nil
}
