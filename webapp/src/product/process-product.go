package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case http.MethodGet:
		if request.PathParameters["productId"] == "" {
			items, err := GetAllItems()
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, err
			}
			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%v", items), StatusCode: 200}, nil
		} else {
			item, err := GetItemById(request.PathParameters["productId"])
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, err
			}
			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%v", item), StatusCode: 200}, nil
		}

	case http.MethodDelete:
		err := DeleteItem(request.PathParameters["productId"])
		if err != nil {
			return events.APIGatewayProxyResponse{Body: "Item Deletion failed.", StatusCode: 502}, err
		}
		return events.APIGatewayProxyResponse{Body: "Item Deleted.", StatusCode: 200}, nil

	case http.MethodPost:
		item, err := AddItem(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%v", item), StatusCode: 502}, err
		}
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%v", item), StatusCode: 200}, nil

	case http.MethodPut:
		item, err := EditItem(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%v", item), StatusCode: 502}, err
		}
		return events.APIGatewayProxyResponse{Body: "Item Edited.", StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{Body: "Not Supported.", StatusCode: 502}, nil
}

func main() {
	//lambda.Start(Handler)
	timeNow := time.Now()
	fmt.Println(timeNow)
}
