package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// basic skeleton for the redirect

// Handler testing to see if this lambda implementation really works

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// logging
	fmt.Printf("event.HTTPMethod %v\n", request.HTTPMethod)
	fmt.Printf("event.Body %v\n", request.Body)
	fmt.Printf("event.QueryStringParameters %v\n", request.QueryStringParameters)
	fmt.Printf("event %v\n", request)

	hash := "super secret code" // placeholder

	// this would be how we get the hash from url
	//if request.HTTPMethod == "GET" {
	//	hash = request.QueryStringParameters["hash"]
	//}
	//

	//  getS3URL(){}
	// 	createTempURL(){}
	//  redirect(){} TODO look into terraform modules might be able to return url to api and dynamically redirect user to URL using terraform

	body := fmt.Sprintf("{\"message\": \"Redirect coming soon\", \"hash\": \"%s\"}", hash)

	// the rest of the error codes are to be handled in terraform, specifically the aws_api_integration_response
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 302, // redirect code
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control":               "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
