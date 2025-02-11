package src

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUserById(id string) (*User, error)
}

type UserRepositoryDynamoDB struct {
	client *dynamodb.Client
	ctx    context.Context
}

func (repo *UserRepositoryDynamoDB) CreateUser(user *User) error {
	return nil
}
