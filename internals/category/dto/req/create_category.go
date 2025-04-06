package dtoreq

type CreateCategoryReq struct {
	Name     string `json:"name" validate:"required,min=3"`
	ParentID *uint  `json:"parentId" validate:"omitempty,numeric"`
}
