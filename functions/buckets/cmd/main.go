//+build linux,amd64,!cgo

package main

import (
	"buckets"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var router *ginadapter.GinLambda

// AWS Lambda Gin Router adapter
func ginAdapter(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return router.ProxyWithContext(ctx, req)
}

//Init all the necessary services
func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	})
	if err != nil {
		panic(err)
	}
	ddb := dynamodb.New(sess)
	r := buckets.NewDynamoDBRepository(ddb)
	s := buckets.NewService(r)
	h := buckets.NewHandler(s)
	router = buckets.NewLambdaRouter(h)
}

func main() {
	lambda.Start(ginAdapter)
}
