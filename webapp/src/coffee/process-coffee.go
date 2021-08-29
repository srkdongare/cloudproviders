package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("Hi from Lambda...")
	stringItem := "Reached Lambda..."

	return events.APIGatewayProxyResponse{Body: stringItem, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
