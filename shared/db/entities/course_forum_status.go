package entities

type CourseForumStatus string

const (
	CourseForumStatus_Open           CourseForumStatus = "open"
	CourseForumStatus_Close          CourseForumStatus = "close"
	CourseForumStatus_CloseTemporary CourseForumStatus = "close-temporary"
)
