package entity

import "slices"

type CourseStatus string

const (
	CourseStatus_Starting   CourseStatus = "starting"
	CourseStatus_InProgress CourseStatus = "in-progress"
	CourseStatus_Done       CourseStatus = "done"
	CourseStatus_Pause      CourseStatus = "pause"
	CourseStatus_Cancel     CourseStatus = "cancel"
)

func (courseStatus CourseStatus) IsValid(canBeEmpty bool) bool {
	if canBeEmpty && courseStatus == "" {
		return true
	}
	courseStatuses := []CourseStatus{
		CourseStatus_Starting,
		CourseStatus_InProgress,
		CourseStatus_Done,
		CourseStatus_Pause,
		CourseStatus_Cancel,
	}

	return slices.Contains(courseStatuses, courseStatus)
}
