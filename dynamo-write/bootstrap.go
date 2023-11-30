package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

type LambdaResponse struct {
	Message string
}

func HandleLambdaEvent(sqsEvent events.SQSEvent) (*LambdaResponse, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config.LoadDefaultConfig failed: %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	for _, record := range sqsEvent.Records {
		if err := handleMessage(ctx, client, record); err != nil {
			fmt.Println(err)
		}
	}
	return &LambdaResponse{Message: "OK"}, nil
}

type sqsSchema struct {
	Link    string
	Title   string
	PubDate time.Time
}

type dynamodbSchema struct {
	Link    string    `dynamodbav:"link"`
	Title   string    `dynamodbav:"title"`
	PubDate time.Time `dynamodbav:"pubdate"`
	TTL     int64
}

func handleMessage(ctx context.Context, dynamoClient *dynamodb.Client, msg events.SQSMessage) error {
	fmt.Printf("%v\n", msg)
	sqsItem := sqsSchema{}
	if err := json.Unmarshal([]byte(msg.Body), &sqsItem); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	ttl := time.Now().Unix() + 24*3600 // 1day
	dynamoItem := dynamodbSchema{Link: sqsItem.Link, Title: sqsItem.Title, PubDate: sqsItem.PubDate, TTL: ttl}
	item, err := attributevalue.MarshalMap(dynamoItem)
	if err != nil {
		return fmt.Errorf("attributevalue.MarshalMap failed: %w", err)
	}
	tableName := os.Getenv("DYNAMODB_TABLE")
	if _, err := dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{Item: item, TableName: aws.String(tableName)}); err != nil {
		return fmt.Errorf("failed to put item on DynamoDB: %w", err)
	}
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
