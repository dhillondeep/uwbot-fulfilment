package main

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"os"
	"strings"
	"warrior_bot/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dhillondeep/go-uw-api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var uwApiClient uwapi.UWAPI

// handles webhook requests from dialogflow using gin
func ginHandler(c *gin.Context) {
	respStr, err := handlers.HandleWebhook(c.Request.Body, &uwApiClient)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}

	// send response back
	c.String(http.StatusOK, respStr)
}

// handles webhook requests from dialogflow using lambda
func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	respStr, err := handlers.HandleWebhook(strings.NewReader(request.Body), &uwApiClient)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       respStr,
			StatusCode: http.StatusBadRequest,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       respStr,
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	// create a client for Uwaterloo API
	uwApiClient = uwapi.Create(os.Getenv("UW_API_KEY"))

	if len(os.Getenv("AWS_LAMBDA_USE")) != 0 {
		lambda.Start(lambdaHandler)
	} else {
		r := gin.Default()
		r.POST("/webhook", ginHandler)

		if err := r.Run(); err != nil {
			logrus.WithError(err).Fatal("Couldn't start server")
		}
	}
}
