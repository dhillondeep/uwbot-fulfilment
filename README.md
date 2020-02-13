![CI](https://github.com/dhillondeep/uwbot-fulfilment/workflows/CI/badge.svg)  ![Deployment](https://github.com/dhillondeep/uwbot-fulfilment/workflows/Deployment/badge.svg)

# UWBot Fulfillment
UWBot is a chatbot that provides University of Waterloo related information in a interactive way! Currently it supports:
* Course prerequisites
* Course sections
* Course section schedule
* Course offerings
* Course enrollments

UWBot uses Dialogflow's NLP engine to interact with users and provides various front ends that users can use. Currently following front-ends are supported:
* Facebook Messenger

This repository contains `fulfillment` code for the bot. It is a webhook service that helps the bot provide dynamic information regarding the University of Waterloo.

## Installation
The project is using GO 1.13 modules!

To build the project:
```bash
go build 
```

To simply run the project:
```bash
go run main.go
```

## Usage
This webhook service can be used and deployed in two ways. The first way is a `persistent http server` and send one is `AWS Lambda`. The usage is controlled by environment variable.

### Persistent HTTP server
The project uses [gin framework](https://github.com/gin-gonic/gin) to create an HTTP server that accepts `POST` requests and sends responses back. No configurations are needed and simply running the project will start the server. It uses port `8080`. Before running, make sure to set up environment variable `UW_API_KEY` which is [API key](https://uwaterloo.ca/api/register) need to access University of Waterloo data.

### AWS Lambda
The project has lambda handler that receives `APIGatewayProxyRequest` and sends `APIGatewayProxyResponse`. When this lambda is paired with `AWS API Gateway`, it creates webhook that can be triggered. Make sure to set up `AWS_LAMBDA_USE` environment variable for the lambda because this signals the app to use Lambda environment. Make sure to also set up environment variable `UW_API_KEY` which is [API key](https://uwaterloo.ca/api/register) need to access University of Waterloo data.
