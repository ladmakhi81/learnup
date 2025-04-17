package dtoreq

type CreateCommentReq struct {
	Content  string `json:"content" validate:"required,min=3"`
	CourseId uint   `json:"-"`
	ParentId *uint  `json:"parentId" validate:"omitempty,numeric"`
}
