package responses

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// creates dialogflow webhook response given some data
func createWebhookResponse(jsonData interface{}) (*dialogflow.WebhookResponse, error) {
	marshalledData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	resp := &dialogflow.WebhookResponse{}
	if err := jsonpb.Unmarshal(bytes.NewBuffer(marshalledData), resp); err != nil {
		return nil, err
	}

	return resp, nil
}
