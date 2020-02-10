package models

type facebookPayload struct {
	TemplateType string           `json:"template_type"`
	Elements     []FbCarouselItem `json:"elements"`
}

type facebookAttachment struct {
	Type    string          `json:"type"`
	Payload facebookPayload `json:"payload,"`
}

type FbButton struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type facebook struct {
	Attachment facebookAttachment `json:"attachment"`
}

type payload struct {
	Facebook facebook `json:"facebook,omitempty"`
}

type FbCarouselItem struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle,omitempty"`
	Buttons  []FbButton `json:"buttons,omitempty"`
}

type FbCarousel struct {
	Payload payload `json:"payload"`
}

// CreateFbCarouselCard creates a single carousel card on messenger
func CreateFbCarouselCard(item FbCarouselItem) *FbCarousel {
	return CreateFbCarousel([]FbCarouselItem{item})
}

// CreateFbCarousel creates a carousel cards on messenger
// This supports many cards and should be used as a lit
func CreateFbCarousel(items []FbCarouselItem) *FbCarousel {
	carousel := &FbCarousel{
		Payload: payload{
			Facebook: facebook{
				Attachment: facebookAttachment{
					Type: "template",
					Payload: facebookPayload{
						TemplateType: "generic",
						Elements:     items,
					},
				},
			},
		},
	}

	return carousel
}
