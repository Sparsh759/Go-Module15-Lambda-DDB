package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"log"
)

type DynamoDBConfig struct {
	Region    string
	TableName string
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() (msg string, err error) {
	var ddb DynamoDBConfig
	ddb.TableName = "ddb-table"
	ddb.Region = "ap-southeast-1"
	ddbSess, err := newDDBClient(&ddb)
	if err != nil {
		fmt.Printf("failed to create connection. Error: %s", err.Error())
		return
	}
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("FirstName"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("PhoneNumber"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PhoneNumber"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("FirstName"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(ddb.TableName),
	}
	_, err = ddbSess.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}
	
	fmt.Println("Created the table", ddb.TableName)
	return
}

func newDDBClient(config *DynamoDBConfig) (dynamodbiface.DynamoDBAPI, error) {
	var awsConfig aws.Config
	awsConfig.Region = aws.String(config.Region)
	awsSession, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(awsSession), nil
}
