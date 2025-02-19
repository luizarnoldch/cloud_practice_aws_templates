package src

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository interface {
	SaveUser(user *User) (*User, error)
	SaveManyUsers(users []User) error
	GetUserById(userId string) (*User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(userId string) error
	UpdateUser(user *User, userId string) (*User, error)
	GetAdultsUsers() ([]User, error)
}

type UserRepositoryDynamoDB struct {
	ctx       context.Context
	client    *dynamodb.Client
	tableName string
}

func NewUserRepositoryDynamoDB(ctx context.Context, client *dynamodb.Client, tableName string) UserRepository {
	return &UserRepositoryDynamoDB{
		ctx:       ctx,
		client:    client,
		tableName: tableName,
	}
}

func (repo *UserRepositoryDynamoDB) SaveManyUsers(users []User) error {
	var writeRequests []types.WriteRequest

	for _, item := range users {
		av, err := attributevalue.MarshalMap(item)
		if err != nil {
			log.Printf("error while marshaling user input %v", err)
			return err
		}

		writeRequest := types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		}

		writeRequests = append(writeRequests, writeRequest)

		if len(writeRequests) == 25 {
			_, err := repo.client.BatchWriteItem(repo.ctx, &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					repo.tableName: writeRequests,
				},
			})
			if err != nil {
				log.Printf("error while writing the batch items %v", err)
				return err
			}
			writeRequests = nil
		}
	}

	if len(writeRequests) > 0 {
		_, err := repo.client.BatchWriteItem(repo.ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				repo.tableName: writeRequests,
			},
		})
		if err != nil {
			log.Printf("error while writing the batch items %v", err)
			return err
		}
	}

	log.Printf("many users saved successfully")
	return nil
}

func (repo *UserRepositoryDynamoDB) GetAdultsUsers() ([]User, error) {
	var response []User

	statement := "SELECT * FROM " + repo.tableName + " WHERE age >= 18"
	output, err := repo.client.ExecuteStatement(repo.ctx, &dynamodb.ExecuteStatementInput{
		Statement: &statement,
	})
	if err != nil {
		log.Printf("error while executing statement: %v", err)
		return []User{}, nil
	}

	if err := attributevalue.UnmarshalListOfMaps(output.Items, &response); err != nil {
		log.Printf("error while deserializing the results: %v", err)
		return []User{}, err
	}

	log.Println("adults users retrieved successfully")
	return response, nil
}

func (repo *UserRepositoryDynamoDB) UpdateUser(user *User, userId string) (*User, error) {
	userMap, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Printf("error while marshaling the user input: %v", err)
		return &User{}, err
	}

	delete(userMap, "ID")

	var updateBuilder expression.UpdateBuilder
	for attr, value := range userMap {
		updateBuilder = updateBuilder.Set(expression.Name(attr), expression.Value(value))
	}

	keycond := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: userId},
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		log.Printf("error while building the UpdateUser expression %v", err)
		return &User{}, err
	}

	input := &dynamodb.UpdateItemInput{
		Key:                       keycond,
		TableName:                 aws.String(repo.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	}

	output, err := repo.client.UpdateItem(repo.ctx, input)
	if err != nil {
		log.Printf("error while updating item %v", err)
		return &User{}, err
	}

	var updatedUser User
	err = attributevalue.UnmarshalMap(output.Attributes, &updatedUser)
	if err != nil {
		log.Printf("error while unmarshaling the response from UpdateUser %v", err)
		return &User{}, err
	}

	return &updatedUser, nil
}

func (repo *UserRepositoryDynamoDB) DeleteUser(userId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: &repo.tableName,
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: userId},
		},
	}

	_, err := repo.client.DeleteItem(repo.ctx, input)
	if err != nil {
		log.Printf("error while deleting user with ID %s: %v", userId, err)
		return err
	}

	log.Printf("user with ID %s deleted successfully", userId)
	return nil
}

func (repo *UserRepositoryDynamoDB) GetAllUsers() ([]User, error) {
	var response []User

	input := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	output, err := repo.client.Scan(repo.ctx, input)
	if err != nil {
		log.Printf("error while executing scan: %v", err)
		return response, err
	}

	if err := attributevalue.UnmarshalListOfMaps(output.Items, &response); err != nil {
		log.Printf("error while deserializing the result: %v", err)
		return nil, err
	}

	log.Println("users retrieved successfully")
	return response, nil
}

func (repo *UserRepositoryDynamoDB) GetUserById(userId string) (*User, error) {
	keycond := expression.Key("ID").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithKeyCondition(keycond).Build()
	if err != nil {
		log.Printf("error while building the GetUserById query expresion %v", err)
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(repo.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	output, err := repo.client.Query(repo.ctx, input)
	if err != nil {
		log.Printf("error while executing the query: %v", err)
		return nil, err
	}

	if len(output.Items) == 0 {
		return &User{}, nil
	}

	var user User
	err = attributevalue.UnmarshalMap(output.Items[0], &user)
	if err != nil {
		log.Printf("error while deserializing the result: %v", err)
		return nil, err
	}

	log.Println("user retrieve successfully")
	return &user, nil
}

func (repo *UserRepositoryDynamoDB) SaveUser(user *User) (*User, error) {
	marshalUser, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Printf("error while marshaling  the user input: %v", err)
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      marshalUser,
		TableName: aws.String(repo.tableName),
	}

	_, err = repo.client.PutItem(repo.ctx, input)
	if err != nil {
		log.Printf("error while saving user on the database: %v", err)
		return nil, err
	}

	log.Println("user created Successfully")
	return user, nil
}
