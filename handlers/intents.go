package handlers

// Course Intents
const (
	CourseTermAvailability = "01_course_term_availability"
	CourseAvailabilityGivenTerm = "02_course_availability_given_term"
	CourseSections = "03_course_section"
	CoursePrerequisites = "04_course_prerequisites"
	CourseSectionSchedule = "05_course_section_schedule"
)

// Term Intents
const (
	CourseAvailNextTerm = "06_term_course_avail_next_term"
	CourseAvailPrevTerm = "07_term_course_avail_prev_term"
	CourseAvailCurrTerm = "08_term_course_avail_curr_term"
	CourseEnrolmentInfo = "09_term_course_enrollment_info"
)

// Intent categories
const (
	CourseIntent = "course"
	TermIntent = "term"
)
