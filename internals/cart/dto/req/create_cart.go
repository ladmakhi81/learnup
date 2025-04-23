package dtoreq

type CreateCartReqDto struct {
	CourseID uint `json:"courseId" validate:"required,gte=1"`
}
