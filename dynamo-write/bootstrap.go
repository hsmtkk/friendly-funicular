package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaRequest struct {
	Link    string
	Title   string
	PubDate time.Time
}

type LambdaResponse struct {
	Message string
}

func HandleLambdaEvent(event *LambdaRequest) (*LambdaResponse, error) {
	fmt.Println(event)
	return &LambdaResponse{Message: "OK"}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
