package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Todo struct {
	ID   string `json:"id"`
	Task string `json:"task"`
	Date string `json:"date"`
}

var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	svc = dynamodb.New(sess)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return createTodo(request)
	case "GET":
		return listTodos(request)
	case "PUT":
		return updateTodo(request)
	case "DELETE":
		return deleteTodo(request)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func createTodo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var todo Todo
	err := json.Unmarshal([]byte(request.Body), &todo)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	item, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("todos"),
		Item:      item,
	})

	if err != nil {
		log.Println("Erro ao salvar item no DynamoDB:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       fmt.Sprintf("Tarefa '%s' criada com sucesso!", todo.Task),
	}, nil
}

func listTodos(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	date := request.PathParameters["date"]
	if date == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String("todos"),
		IndexName:              aws.String("DateIndex"),
		KeyConditionExpression: aws.String("#date = :date"),
		ExpressionAttributeNames: map[string]*string{
			"#date": aws.String("date"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":date": {S: aws.String(date)},
		},
	}

	result, err := svc.Query(input)
	if err != nil {
		log.Println("Erro ao buscar tarefas:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	var todos []Todo
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	body, _ := json.Marshal(todos)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil
}

func updateTodo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	var todo Todo
	err := json.Unmarshal([]byte(request.Body), &todo)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	item, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("todos"),
		Item:      item,
	})

	if err != nil {
		log.Println("Erro ao atualizar tarefa:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Tarefa atualizada com sucesso"}, nil
}

func deleteTodo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if id == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("todos"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})

	if err != nil {
		log.Println("Erro ao excluir tarefa:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "Tarefa excluída com sucesso"}, nil
}

func main() {
	lambda.Start(handler)
}
