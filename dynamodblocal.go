package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateLocalClient() *dynamodb.DynamoDB {

	creds := credentials.NewStaticCredentials("123", "123", "")
	awsConfig := &aws.Config{
		Credentials: creds,
	}
	awsConfig.WithRegion("us-east-1")
	awsConfig.WithEndpoint("http://localhost:8000")

	s, err := session.NewSession(awsConfig)
	if err != nil {
		panic(err)
	}
	dynamodbconn := dynamodb.New(s)
	return dynamodbconn
}

/*
func tableExists(d *dynamodb.Client, name string) bool {
	tables, err := d.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}
	for _, n := range tables.TableNames {
		if n == name {
			return true
		}
	}
	return false
}
*/
