package models

type FbButton struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type facebookPayload struct {
	TemplateType string           `json:"template_type"`
	Elements     []FbCarouselItem `json:"elements"`
}

type facebookAttachment struct {
	Type    string          `json:"type"`
	Payload facebookPayload `json:"payload,"`
}

type facebook struct {
	Attachment facebookAttachment `json:"attachment"`
}

type FbCarouselItem struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle,omitempty"`
	Buttons  []FbButton `json:"buttons,omitempty"`
}

// CreateFbCarouselCard creates a single carousel card on messenger
func CreateFbCarouselCard(item FbCarouselItem) *DialogflowResponse {
	return CreateFbCarousel([]FbCarouselItem{item})
}

// CreateFbCarousel creates a carousel cards on messenger
// This supports many cards and should be used as a lit
func CreateFbCarousel(items []FbCarouselItem) *DialogflowResponse {
	itemsShow := items

	if len(itemsShow) > 10 {
		itemsShow = itemsShow[:10]
	}

	carousel := &DialogflowResponse{
		Payload: payload{
			Facebook: facebook{
				Attachment: facebookAttachment{
					Type: "template",
					Payload: facebookPayload{
						TemplateType: "generic",
						Elements:     itemsShow,
					},
				},
			},
		},
	}

	return carousel
}
