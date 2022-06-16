package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func main() {
	fmt.Println("Start")
	getRandomProject()
	fmt.Println("Close")
}

func getRandomProject() {
	d := CreateLocalClient()

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableProject),
	}

	out, err := d.Scan(input)
	if err != nil {
		log.Fatal("scan failed", err)
	}
	fmt.Println(out.Items)

	projects, err := toProjectItems(out.Items)

	fmt.Println(projects)
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
