package models

import (
	"github.com/dhillondeep/go-uw-api"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Fields struct {
	Subject    string
	CatalogNum string
	Term       string
	Section    string
}

type RespContext struct {
	StatusCode int
	Resp       *DialogflowResponse
}

type ReqContext struct {
	UWApiClient       *uwapi.UWAPI
	DialogflowRequest *dialogflow.WebhookRequest
	Fields            *Fields
}

type Payload struct {
	Facebook *Facebook `json:"facebook,omitempty"`
}

type DialogflowResponse struct {
	FulfillmentText string   `json:"fulfillment_text,omitempty"`
	Payload         *Payload `json:"payload,omitempty"`
	Source          string   `json:"source,omitempty"`
}
