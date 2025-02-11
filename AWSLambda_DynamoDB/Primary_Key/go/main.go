package main

import (
	"context"
	"log"
	"main/lib"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	DYNAMODB_LOCAL_TABLE = "Users"
)

func CreateTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	},
	)
	if err != nil {
		return err
	}
	log.Println("Table created successfully", tableName)
	return nil
}

func DeleteTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return err
	}

	log.Println("Table deleted successfully", tableName)

	return nil
}

func ListTables(ctx context.Context, client *dynamodb.Client) {
	result, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalln("Error while listting tables", err)
	}

	for _, table := range result.TableNames {
		log.Println(" -", table)
	}
}

func main() {
	ctx := context.TODO()

	client, err := lib.NewLocalDynamoDBClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client %v", err)
	}

	DeleteTable(ctx, client, DYNAMODB_LOCAL_TABLE)

	// err = CreateTable(ctx, client, DYNAMODB_LOCAL_TABLE)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// ListTables(ctx, client)
}
