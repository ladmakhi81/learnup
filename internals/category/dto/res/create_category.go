package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
	"time"
)

type CreateCategoryResDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCreateCategoryResDto(category *entities.Category) CreateCategoryResDto {
	return CreateCategoryResDto{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}
}
