package dtoreq

type CreateCartReq struct {
	CourseID uint `json:"courseId" validate:"required,gte=1"`
}
