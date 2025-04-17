package dtoreq

type CreateCommentReq struct {
	Content  string `json:"content" validate:"required,min=3"`
	CourseId uint   `json:"courseId" validate:"required,gte=1"`
	ParentId *uint  `json:"parentId" validate:"omitempty,numeric"`
}
