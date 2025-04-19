package dtoreq

type CreateCartReq struct {
	UserID   uint `json:"-"`
	CourseID uint `json:"courseId" validate:"required,gte=1"`
}
