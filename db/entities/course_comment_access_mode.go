package entities

import "slices"

type CourseCommentAccessMode string

const (
	CourseCommentAccessMode_All      = "all"
	CourseCommentAccessMode_Students = "students"
)

func (courseCommentAccessMode CourseCommentAccessMode) IsValid(canBeEmpty bool) bool {
	if courseCommentAccessMode == "" && canBeEmpty {
		return true
	}
	courseCommentAccessModes := []CourseCommentAccessMode{
		CourseCommentAccessMode_Students,
		CourseCommentAccessMode_All,
	}
	return slices.Contains(courseCommentAccessModes, courseCommentAccessMode)
}
