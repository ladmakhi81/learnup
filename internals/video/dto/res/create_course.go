package dtores

type CreateCourseRes struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	CourseID *uint  `json:"courseId"`
}

func NewCreateCourseRes(id uint, url string, courseID *uint) CreateCourseRes {
	return CreateCourseRes{
		ID:       id,
		URL:      url,
		CourseID: courseID,
	}
}
