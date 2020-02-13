package responses

import (
	"fmt"
	"net/http"
	"strings"
	"uwbot/models"
)

// TextResponse creates dialogflow webhook response for simple text
func TextResponse(text string) *models.RespContext {
	return &models.RespContext{
		StatusCode: http.StatusOK,
		Resp: &models.DialogflowResponse{
			FulfillmentText: strings.TrimSpace(text),
		},
	}
}

// TextResponsef creates dialogflow webhook response for simple text but,
// it allows the text to be formatted
func TextResponsef(format string, a ...interface{}) *models.RespContext {
	return TextResponse(fmt.Sprintf(format, a...))
}

func FbCarouselCard(item *models.FbCarouselItem) *models.RespContext {
	return FbCarousel([]*models.FbCarouselItem{item})
}

func FbCarousel(items []*models.FbCarouselItem) *models.RespContext {
	itemsShow := items

	// facebook limitation and does not render if num cards > 10
	if len(itemsShow) > 10 {
		itemsShow = itemsShow[:10]
	}

	return &models.RespContext{
		StatusCode: http.StatusOK,
		Resp: &models.DialogflowResponse{
			Payload: &models.Payload{
				Facebook: &models.Facebook{
					Attachment: models.FbAttachment{
						Type: "template",
						Payload: models.FbPayload{
							TemplateType: "generic",
							Elements:     itemsShow,
						},
					},
				},
			},
		},
	}
}
