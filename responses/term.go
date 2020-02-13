package responses

import (
	"fmt"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
)

func CourseNotOfferedNextTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is not offered next term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedNextTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is offered next term!", fields.Subject, fields.CatalogNum)
}

func CourseNotOfferedPrevTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s was not offered last term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedPrevTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s was offered last term!", fields.Subject, fields.CatalogNum)
}

func CourseNotOfferedCurrTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is not offered this term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedCurrTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is being offered this term!", fields.Subject, fields.CatalogNum)
}

func SectionEnrollmentInfoItem(enrollmentInfo string, fields *models.Fields) *models.FbCarouselItem {
	return &models.FbCarouselItem{
		Title:    fmt.Sprintf("%s %s %s", fields.Subject, fields.CatalogNum, fields.Section),
		Subtitle: strings.TrimSpace(enrollmentInfo),
		Buttons: []models.FbButton{
			{
				Type:  "web_url",
				Url:   fmt.Sprintf(helpers.UWFlowCourseUrl, fields.Subject, fields.CatalogNum),
				Title: "More Info",
			},
		},
	}
}
