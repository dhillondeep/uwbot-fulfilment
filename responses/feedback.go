package responses

import (
	"uwbot/models"
)

/********* Course Feedback *********/

func CourseNotFound(subject, catalogNum string) *models.RespContext {
	return TextResponsef("Sorry, I am unable to find %s %s", subject, catalogNum)
}

func CourseOfferingNotFound(subject, catalogNum string) *models.RespContext {
	return TextResponsef("There are no offerings for %s %s", subject, catalogNum)
}

func CoursePrereqNotFound(subject, catalogNum string) *models.RespContext {
	return TextResponsef("There are no prerequisites for %s %s", subject, catalogNum)
}

/********* Term Feedback *********/

func NoCourseSecFound(subject, catalogNum string) *models.RespContext {
	return TextResponsef("Sorry, I am unable to find any section for %s %s", subject, catalogNum)
}
