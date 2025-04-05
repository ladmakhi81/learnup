package dtoreq

type LoginReq struct {
	Phone    string `json:"phone" validate:"required,len=11"`
	Password string `json:"password" validate:"required,min=8"`
}
