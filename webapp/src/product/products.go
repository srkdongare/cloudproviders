package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func GetItemById(id string) (Item, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	var item Item

	// Perform the query
	fmt.Println("Trying to read from table: ", os.Getenv("TABLE_NAME"))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	return item, nil
}

func GetAllItems() (Item, error) {

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	item := Item{}

	// Perform the query
	fmt.Println("Trying to read from table: ", os.Getenv("TABLE_NAME"))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	return item, nil
}

func AddItem(body string) (Item, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var thisItem Item
	json.Unmarshal([]byte(body), &thisItem)

	thisItem.Id = time.Now().UnixNano()

	fmt.Println("Item to add:", thisItem)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(thisItem)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return thisItem, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	return thisItem, err

}

func DeleteItem(id string) error {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func EditItem(body string) (Item, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var thisItem Item
	json.Unmarshal([]byte(body), &thisItem)

	fmt.Println("Item to update:", thisItem)

	// Update Item in table and return
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(thisItem.Name),
			},
			":description": {
				S: aws.String(thisItem.Description),
			},
			":price": {
				N: aws.String(strconv.FormatFloat(thisItem.Price, 'f', 2, 64)),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.FormatInt(thisItem.Id, 10)),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set info.rating = :r, info.plot = :p"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return thisItem, err

}
