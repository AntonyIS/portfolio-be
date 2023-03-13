/*
Package name : repository
File name : dynamodb.go
Author : Antony Injila
Description :
	- Host dynamoDb database specific methods
*/

package repostitory

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/joho/godotenv"
)

type dynamoDbClient struct {
	client            *dynamodb.DynamoDB
	usersTableName    string
	projectsTableName string
}

func NewDynaDBRepository() ports.PortfolioRepository {
	// Load portifolio environmental variables
	loadEnv()
	// Add AWS session for DynamoDB
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		},
	}))

	var (
		userTablename    = os.Getenv("DYNAMODB_USERS_TABLE")
		projectTablename = os.Getenv("DYNAMODB_PROJECTS_TABLE")
	)
	return &dynamoDbClient{
		client:            dynamodb.New(sess),
		usersTableName:    userTablename,
		projectsTableName: projectTablename,
	}
}

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}

func (db *dynamoDbClient) CreateUser(user *domain.User) (*domain.User, error) {

	entityParsed, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.usersTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, err
	}

	user, err = db.ReadUser(user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *dynamoDbClient) ReadUser(id string) (*domain.User, error) {
	result, err := db.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New(fmt.Sprintf("user with id [ %s ] not found", id))
	}
	var user domain.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *dynamoDbClient) ReadUsers() ([]*domain.User, error) {
	users := []*domain.User{}
	filt := expression.Name("Id").AttributeNotExists()
	proj := expression.NamesList(
		expression.Name("id"),
		expression.Name("firstname"),
		expression.Name("lastname"),
		expression.Name("email"),
		expression.Name("projects"),
	)
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.usersTableName),
	}
	result, err := db.client.Scan(params)

	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		var user domain.User

		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}

	return users, nil
}

func (db *dynamoDbClient) UpdateUser(user *domain.User) (*domain.User, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.usersTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *dynamoDbClient) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.usersTableName),
	}

	res, err := db.client.DeleteItem(input)
	if res == nil {
		return errors.New(fmt.Sprintf("No user to delete: %v", err))
	}
	if err != nil {
		return errors.New(fmt.Sprintf("Got error calling DeleteItem: %v", err))
	}
	return nil
}

func (db *dynamoDbClient) CreateProject(project *domain.Project) (*domain.Project, error) {

	entityParsed, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.projectsTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (db *dynamoDbClient) ReadProject(id string) (*domain.Project, error) {
	result, err := db.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.projectsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New(fmt.Sprintf("project with id [ %s ] not found", id))
	}
	var project domain.Project
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (db *dynamoDbClient) ReadProjects() ([]*domain.Project, error) {
	projects := []*domain.Project{}
	filt := expression.Name("Id").AttributeNotExists()
	proj := expression.NamesList(
		expression.Name("id"),
		expression.Name("title"),
		expression.Name("body"),
		expression.Name("user_id"),
		expression.Name("created_at"),
	)
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.projectsTableName),
	}
	result, err := db.client.Scan(params)
	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		var project domain.Project

		err = dynamodbattribute.UnmarshalMap(item, &project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)

	}

	return projects, nil
}

func (db *dynamoDbClient) UpdateProject(project *domain.Project) (*domain.Project, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.projectsTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (db *dynamoDbClient) DeleteProject(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.projectsTableName),
	}

	res, err := db.client.DeleteItem(input)
	if res == nil {
		return errors.New(fmt.Sprintf("no project to delete: %s", err))
	}
	if err != nil {
		return errors.New(fmt.Sprintf("got error calling DeleteItem: %s", err))
	}
	return nil
}
