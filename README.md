# UWBot Fulfilment
Dialogflow fulfilment code for UWBot. This code is called by webhook and response is returned
by providing dynamic data.

## Installation
The project is using go 1.13 modules!

To build the project do:
```bash
go build 
```

To simply run the project do:
```bash
go run main.go
```

## Usage
There are two environments that are supported by the project. While developing, local env can be used
which creates a server. For deployment, lambda env is used where lambda handler is used.

### Local
For local environment, `gin` server is started at port `8080`. The webhook is accessible through `/webhook`
endpoint. To use this locally, set up env variable `UW_API_KEY` which is API key for University of Waterloo
Open Data API.

### Lambda
For lambda environment, lambda handler is used. To use this, set up env variable `UW_API_KEY` which is API key 
for University of Waterloo Open Data API. You will also need to set up env variable `AWS_LAMBDA_USE` which
signals that project should use lambda env.
