package src_test

import (
	"context"
	"main/crud_dynamodb/src"
	"main/lib"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserDynamoDBSuite struct {
	suite.Suite
	// Repository
	repository src.UserRepository

	// Setup variables
	client    *dynamodb.Client
	ctx       context.Context
	tableName string

	// Global variables
	userUuid string
}

func (suite *UserDynamoDBSuite) SetupSuite() {
	var err error
	suite.ctx = context.TODO()
	suite.client, err = lib.NewLocalDynamoDBClient(suite.ctx)
	suite.NoError(err, "error while creating dynamodb local client")
	suite.tableName = "UserTable"

	err = lib.CreateDynamoDBUseTable(suite.ctx, suite.client, suite.tableName)
	suite.NoError(err, "error while creating user table")

	suite.repository = src.NewUserRepositoryDynamoDB(suite.ctx, suite.client, suite.tableName)
}

func (suite *UserDynamoDBSuite) TearDownSuite() {
	err := lib.DeleteLocalDynamoDBTable(suite.ctx, suite.client, suite.tableName)
	suite.NoError(err, "error while deleting the user table")
}

func (suite *UserDynamoDBSuite) TestCreateUser() {
	suite.userUuid = uuid.NewString()
	user := src.User{
		Id:   suite.userUuid,
		Name: "newUser",
		Age:  18,
	}

	newUser, err := suite.repository.SaveUser(&user)
	suite.NoError(err, "error while creating new user")

	suite.Equal(user.Id, newUser.Id, "uuid are not equal")
	suite.Equal(user.Name, newUser.Name, "name are not equal")
	suite.Equal(user.Age, newUser.Age, "age are not equal")
}

func (suite *UserDynamoDBSuite) TestGetUserByID() {
	user, err := suite.repository.GetUserById(suite.userUuid)
	suite.NoError(err, "error while getting user by ID")

	suite.Equal(user.Id, suite.userUuid, "uuid are not equal")
}

func (suite *UserDynamoDBSuite) TestGetAllUsers() {
	users, err := suite.repository.GetAllUsers()
	suite.NoError(err, "error while getting all users")

	suite.Equal(len(users), 1, "not the same users on database that injected")
}

func (suite *UserDynamoDBSuite) TestUpdateUser() {
	updateData := &src.User{
		Name: "updateUser",
		Age:  25,
	}

	updatedUser, err := suite.repository.UpdateUser(updateData, suite.userUuid)
	suite.NoError(err, "error while updating user")

	suite.Equal(updateData.Name, updatedUser.Name, "mock and new names are not equal after update")
	suite.Equal(updateData.Age, updatedUser.Age, "mock and new ages are not equal after update")
}

func (suite *UserDynamoDBSuite) TestDeleteUser() {
	userForDeleteUuid := uuid.NewString()
	user := src.User{
		Id:   userForDeleteUuid,
		Name: "newName",
		Age:  17,
	}

	_, err := suite.repository.SaveUser(&user)
	suite.NoError(err, "error while creating new user for deleting")

	usersBefore, err := suite.repository.GetAllUsers()
	suite.NoError(err, "error while getting all users after delete")

	err = suite.repository.DeleteUser(userForDeleteUuid)
	suite.NoError(err, "error while deleting user")

	usersAfter, err := suite.repository.GetAllUsers()
	suite.NoError(err, "error while getting all users after delete")

	suite.Equal(len(usersBefore)-1, len(usersAfter), "user still exists in the database after deletion")
}

func (suite *UserDynamoDBSuite) TestGetAdultUsers() {
	adults, err := suite.repository.GetAdultsUsers()
	suite.NoError(err, "error while getting adults from users")

	suite.Equal(len(adults), 1, "not the same adults on database than expected")
}

func (suite *UserDynamoDBSuite) TestSaveManyUsers() {
	users := []src.User{
		{Id: uuid.NewString(), Name: "John Doe", Age: 16},
		{Id: uuid.NewString(), Name: "Jane Smith", Age: 15},
	}

	err := suite.repository.SaveManyUsers(users)
	suite.NoError(err, "error while saving many users")
}

func TestUserDynamoDBSuite(t *testing.T) {
	suite.Run(t, new(UserDynamoDBSuite))
}
