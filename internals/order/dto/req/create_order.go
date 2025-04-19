package dtoreq

type CreateOrderReq struct {
	UserID uint   `json:"-"`
	Carts  []uint `json:"carts" validate:"required,dive,gte=1"`
}
