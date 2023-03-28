/*
Package name : repository
File name : dynamodb.go
Author : Antony Injila
Description :
	- Host dynamoDb database specific methods
*/

package repository

import (
	"errors"
	"fmt"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	errs "github.com/pkg/errors"
)

var (
	internalServerError = config.ErrInternalServer.Error()
	itemNotFound        = config.ErrNotFound.Error()
	invalidItem         = config.ErrInvalidItem.Error()
)

type dynamoDbClient struct {
	client            *dynamodb.DynamoDB
	usersTableName    string
	projectsTableName string
}

func NewDynamoDBRepository(c *config.AppConfig) ports.PortfolioRepository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	}))
	return &dynamoDbClient{
		client:            dynamodb.New(sess),
		usersTableName:    c.UsersTable,
		projectsTableName: c.ProjectTable,
	}
}

func (db *dynamoDbClient) CreateUser(user *domain.User) (*domain.User, error) {

	entityParsed, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateUser")
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.usersTableName),
	}

	_, err = db.client.PutItem(input)

	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateUser")
	}

	user, err = db.ReadUserWithEmail(user.Email)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateUser")
	}

	return user, nil
}
func (db *dynamoDbClient) ReadUser(id string) (*domain.User, error) {
	users, err := db.ReadUsers()
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadUser")
	}
	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", itemNotFound, err)), "adapters.repository.dynamodb.ReadUser")
}
func (db *dynamoDbClient) ReadUserWithEmail(email string) (*domain.User, error) {
	result, err := db.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadUserWithEmail")
	}

	if result.Item == nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", itemNotFound, err)), "adapters.repository.dynamodb.ReadUserWithEmail")
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
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadUsers")
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
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadUsers")
	}

	for _, item := range result.Items {
		var user domain.User

		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadUsers")
		}

		users = append(users, &user)

	}
	return users, nil
}
func (db *dynamoDbClient) UpdateUser(user *domain.User) (*domain.User, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.UpdateUser")
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.usersTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.UpdateUser")
	}

	return user, nil
}

func (db *dynamoDbClient) DeleteUser(email string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(db.usersTableName),
	}

	res, err := db.client.DeleteItem(input)
	if res == nil {
		return errs.Wrap(errors.New(fmt.Sprintf("%s: %s", itemNotFound, err)), "adapters.repository.dynamodb.DeleteUser")
	}
	if err != nil {
		return errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.DeleteUser")
	}
	return nil
}

func (db *dynamoDbClient) CreateProject(project *domain.Project) (*domain.Project, error) {

	entityParsed, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateProject")
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.projectsTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateProject")
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
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProject")
	}
	if result.Item == nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProject")
	}
	var project domain.Project
	err = dynamodbattribute.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProject")
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
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProjects")
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
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProjects")
	}

	for _, item := range result.Items {
		var project domain.Project

		err = dynamodbattribute.UnmarshalMap(item, &project)
		if err != nil {
			return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.ReadProjects")
		}
		projects = append(projects, &project)

	}

	return projects, nil
}

func (db *dynamoDbClient) UpdateProject(project *domain.Project) (*domain.Project, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.UpdateProject")
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.projectsTableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.UpdateProject")
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
		return errs.Wrap(errors.New(fmt.Sprintf("%s: %s", itemNotFound, err)), "adapters.repository.dynamodb.DeleteProject")
	}
	if err != nil {
		return errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.DeleteProject")
	}
	return nil
}
