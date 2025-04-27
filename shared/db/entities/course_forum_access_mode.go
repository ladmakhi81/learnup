package entities

type CourseForumAccessMode string

const (
	CourseForumAccessMode_Student        CourseForumAccessMode = "student-only"
	CourseForumAccessMode_Teacher        CourseForumAccessMode = "teacher-only"
	CourseForumAccessMode_StudentTeacher CourseForumAccessMode = "student-teacher"
)
