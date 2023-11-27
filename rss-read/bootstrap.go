package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	sess := session.Must(session.NewSession())
	sqsClient := sqs.New(sess)
	for _, feed := range feeds {
		msgID, err := handleFeed(sqsClient, queueURL, feed)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("SQS Message ID: %s\n", msgID)
	}
	return &LambdaResponse{Message: "OK"}, nil
}

func handleFeed(sqsClient *sqs.SQS, queueURL string, feed rss.Feed) (string, error) {
	encoded, err := json.Marshal(feed)
	if err != nil {
		return "", fmt.Errorf("json.Marshal failed: %w", err)
	}
	msg := string(encoded)
	resp, err := sqsClient.SendMessage(&sqs.SendMessageInput{MessageBody: &msg, QueueUrl: &queueURL})
	if err != nil {
		return "", fmt.Errorf("sendMessage failed: %w", err)
	}
	return *resp.MessageId, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
