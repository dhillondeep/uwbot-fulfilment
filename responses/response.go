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
		Resp:       models.CreateTextResponse(strings.Trim(strings.TrimSpace(text), ",")),
	}
}

// TextResponsef creates dialogflow webhook response for simple text but,
// it allows the text to be formatted
func TextResponsef(format string, a ...interface{}) *models.RespContext {
	return TextResponse(fmt.Sprintf(format, a...))
}

func FbCarouselCard(item *models.FbCarouselItem) *models.RespContext {
	return &models.RespContext{
		StatusCode: http.StatusOK,
		Resp:       models.CreateFbCarouselCard(item),
	}
}

func FbCarousel(items []*models.FbCarouselItem) *models.RespContext {
	return &models.RespContext{
		StatusCode: http.StatusOK,
		Resp:       models.CreateFbCarousel(items),
	}
}
