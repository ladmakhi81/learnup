package dtores

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"time"
)

type CategoryPageableItemDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parentId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func MapCategoryPageableItemsDto(categories []*entities.Category) []*CategoryPageableItemDto {
	pageableItems := make([]*CategoryPageableItemDto, len(categories))
	for categoryIndex, category := range categories {
		pageableItems[categoryIndex] = &CategoryPageableItemDto{
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			DeletedAt: category.DeletedAt.Time,
			ID:        category.ID,
			ParentID:  category.ParentCategoryID,
		}
	}
	return pageableItems
}
