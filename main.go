package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, req events.CloudWatchEvent) {
	log.Printf("Handling new request")
}

func main() {
	log.Printf("New execution context created")
	lambda.Start(handleRequest)
}
