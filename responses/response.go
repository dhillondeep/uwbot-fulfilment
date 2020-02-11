package responses

import (
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"strings"
	"warrior_bot/models"
)

// TextResponse creates dialogflow webhook response for simple text
func TextResponse(text string) (*dialogflow.WebhookResponse, error) {
	return createWebhookResponse(models.CreateTextResponse(strings.Trim(strings.TrimSpace(text), ",")))
}

// TextResponsef creates dialogflow webhook response for simple text but,
// it allows the text to be formatted
func TextResponsef(format string, a ...interface{}) (*dialogflow.WebhookResponse, error) {
	return TextResponse(fmt.Sprintf(format, a...))
}

func FbCarouselCard(item models.FbCarouselItem) (*dialogflow.WebhookResponse, error) {
	return createWebhookResponse(models.CreateFbCarouselCard(item))
}

func FbCarousel(items []models.FbCarouselItem) (*dialogflow.WebhookResponse, error) {
	return createWebhookResponse(models.CreateFbCarousel(items))
}
