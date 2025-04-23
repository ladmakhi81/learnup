package entities

type NotificationType string

const (
	NotificationType_CompleteVideoUpload                   = "complete-video-upload"
	NotificationType_CompleteIntroductionCourseVideoUpload = "complete-introduction-course-video-upload"
	NotificationType_CourseVerified                        = "course-verified"
)
