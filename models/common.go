package models

type DialogflowResponse struct {
	FulfillmentText string `json:"fulfillment_text"`
}

// CreateTextResponse creates basic text response
func CreateTextResponse(text string) *DialogflowResponse {
	return &DialogflowResponse{
		FulfillmentText: text,
	}
}
