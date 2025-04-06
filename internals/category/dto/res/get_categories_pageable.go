package dtores

import (
	"github.com/ladmakhi81/learnup/internals/category/entity"
	"time"
)

type CategoryPageableItem struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parentId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func MapCategoriesToPageableItems(categories []*entity.Category) []*CategoryPageableItem {
	pageableItems := make([]*CategoryPageableItem, len(categories))
	for categoryIndex, category := range categories {
		pageableItems[categoryIndex] = &CategoryPageableItem{
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
