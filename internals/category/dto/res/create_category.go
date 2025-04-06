package dtores

import "time"

type CreateCategoryRes struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCreateCategoryRes(
	id uint,
	name string,
	createdAt time.Time,
) CreateCategoryRes {
	return CreateCategoryRes{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
	}
}
