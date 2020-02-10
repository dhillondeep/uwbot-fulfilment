package responses

import (
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

/********* Course Errors *********/

func CourseNotFoundError(subject, catalogNum string) (*dialogflow.WebhookResponse, error) {
	respStr := "Sorry, I was unable to find this information about %s %s at University of Waterloo"
	return TextResponsef(respStr, subject, catalogNum)
}

func CourseOfferingError(subject, catalogNum string) (*dialogflow.WebhookResponse, error) {
	respStr := "Sorry, I was unable to find any offerings for %s %s"
	return TextResponsef(respStr, subject, catalogNum)
}

func CoursePrerequisitesNotFound(subject, catalogNum string) (*dialogflow.WebhookResponse, error) {
	respStr := "Sorry, I was unable to find information about prerequisites for %s %s"
	return TextResponsef(respStr, subject, catalogNum)
}

func CourseSectionInfoNotFound (subject, catalogNum, section string) (*dialogflow.WebhookResponse, error) {
	respStr := "Sorry, I was unable to find information about %s %s %s"
	return TextResponsef(respStr, subject, catalogNum, section)
}

/********* Term Errors *********/

func CourseNotFoundNextTerm(subject, catalogNum string) (*dialogflow.WebhookResponse, error) {
	respStr := "Sorry, I was "
	return TextResponsef(respStr, subject, catalogNum)
}