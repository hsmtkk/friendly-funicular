package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hsmtkk/friendly-funicular/rss-read/rss"
)

type LambdaRequest struct {
}

type LambdaResponse struct {
	Message string
}

func HandleLambdaEvent(event *LambdaRequest) (*LambdaResponse, error) {
	newsURL := os.Getenv("NEWS_URL")
	queueURL := os.Getenv("QUEUE_URL")
	feeds, err := rss.GetFeeds(newsURL)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config.LoadDefaultConfig failed: %w", err)
	}
	client := sqs.NewFromConfig(cfg)
	for _, feed := range feeds {
		msgID, err := handleFeed(ctx, client, queueURL, feed)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("SQS Message ID: %s\n", msgID)
	}
	return &LambdaResponse{Message: "OK"}, nil
}

func handleFeed(ctx context.Context, sqsClient *sqs.Client, queueURL string, feed rss.Feed) (string, error) {
	encoded, err := json.Marshal(feed)
	if err != nil {
		return "", fmt.Errorf("json.Marshal failed: %w", err)
	}
	msg := string(encoded)
	resp, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: aws.String(msg), QueueUrl: aws.String(queueURL)})
	if err != nil {
		return "", fmt.Errorf("sendMessage failed: %w", err)
	}
	return *resp.MessageId, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
