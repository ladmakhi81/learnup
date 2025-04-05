package dtoreq

type CreateBasicUserReq struct {
	Phone     string `json:"phone" validate:"required,numeric,len=11"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required,min=3"`
	LastName  string `json:"lastName" validate:"required,min=3"`
}
