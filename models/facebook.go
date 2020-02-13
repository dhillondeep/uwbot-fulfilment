package models

type FbButton struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type FbPayload struct {
	TemplateType string            `json:"template_type"`
	Elements     []*FbCarouselItem `json:"elements"`
}

type FbAttachment struct {
	Type    string    `json:"type,omitempty"`
	Payload FbPayload `json:"payload,omitempty"`
}

type Facebook struct {
	Attachment FbAttachment `json:"attachment,omitempty"`
}

type FbCarouselItem struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle"`
	Buttons  []FbButton `json:"buttons,omitempty"`
}
