package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const tableProject = "projects-dev"

type Project struct {
	UUID        string `json:"uuid"`
	ContestUUID string `json:"contest_uuid"`
	UserUUID    string `json:"user_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ProjectUUID struct {
	UUID string `json:"uuid"`
}

var db *dynamodb.DynamoDB

func init() {
	db = CreateLocalClient()
}

func main() {
	fmt.Println("Start")
	getRandomProject()
	fmt.Println("Close")
}

func getRandomProject() {
	fmt.Println("Load UUIDs")

	projectUUIDs, err := getUUIDs()
	if err != nil {
		fmt.Println("Error loading project uuids : ", err)
	}
	fmt.Println(projectUUIDs)

	fmt.Println("Random UUID")
	randomUUID, err := getRandomUUID(projectUUIDs)
	if err != nil {
		fmt.Println("Error get random UUID ", err)
	}
	fmt.Println(randomUUID)

	fmt.Println("Load Project Random")
	getProjectItem(randomUUID)
}

func getUUIDs() ([]ProjectUUID, error) {

	proj := expression.NamesList(expression.Name("uuid"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                aws.String(tableProject),
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	out, err := db.Scan(input)
	if err != nil {
		return nil, err
	}

	projects, err := toProjectUUIDItems(out.Items)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func getRandomUUID(projectUUIDs []ProjectUUID) (ProjectUUID, error) {
	quantity := len(projectUUIDs)
	max := int64(quantity - 1)

	if max == -1 {
		return ProjectUUID{}, errors.New("There are no projects")
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return ProjectUUID{}, nil
	}
	n := nBig.Int64()
	fmt.Println("Random Index: ", n)

	return projectUUIDs[n], nil
}

func getProjectItem(projectUUID ProjectUUID) {
	condition := expression.Name("uuid").Equal(expression.Value(projectUUID.UUID))

	expr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableProject),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Condition(),
	}

	out, err := db.Scan(input)
	if err != nil {
		log.Fatal("scan failed", err)
	}

	projects, err := toProjectItems(out.Items)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(projects)
}

func toProjectUUIDItems(rawItems []map[string]*dynamodb.AttributeValue) ([]ProjectUUID, error) {
	if len(rawItems) == 0 {
		return []ProjectUUID{}, nil
	}

	dbItems := []ProjectUUID{}

	if err := dynamodbattribute.UnmarshalListOfMaps(rawItems, &dbItems); err != nil {
		return nil, err
	}

	return dbItems, nil
}

func toProjectItems(rawItems []map[string]*dynamodb.AttributeValue) ([]Project, error) {
	if len(rawItems) == 0 {
		return []Project{}, nil
	}

	dbItems := []Project{}

	if err := dynamodbattribute.UnmarshalListOfMaps(rawItems, &dbItems); err != nil {
		return nil, err
	}

	return dbItems, nil
}
