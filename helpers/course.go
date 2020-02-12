package helpers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"strings"
)

func CreateCourseSectionSchedule(class *gabs.Container) string {
	startDate := class.Path("date.start_date").Data()
	endDate := class.Path("date.end_date").Data()

	startTime := class.Path("date.start_time").Data()
	endTime := class.Path("date.end_time").Data()
	weekdays := class.Path("date.weekdays").Data()

	building := class.Path("location.building").Data()
	room := class.Path("location.room").Data()

	var schStrClass string

	if startDate != nil {
		schStrClass += fmt.Sprintf("Start Date: %s\n", startDate.(string))
	}

	if endDate != nil {
		schStrClass += fmt.Sprintf("End Date: %s\n", endDate.(string))
	}

	if startTime != nil && endTime != nil && weekdays != nil {
		schStrClass += fmt.Sprintf("Timmings: %s - %s %s\n",
			startTime.(string), endTime.(string), weekdays.(string))
	}

	if building != nil && room != nil {
		schStrClass += fmt.Sprintf("Building: %s %s\n", building.(string), room.(string))
	}

	return schStrClass
}

func CleanPrereqString(prerequisites string) string {
	return strings.Trim(strings.Replace(prerequisites, "Prereq: ", "", 1), ".")
}
