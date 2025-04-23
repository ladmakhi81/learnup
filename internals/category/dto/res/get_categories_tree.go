package dtores

import (
	"github.com/ladmakhi81/learnup/shared/db/entities"
)

type GetCategoriesTreeItemDto struct {
	ID            uint                        `json:"id"`
	Name          string                      `json:"name"`
	ParentID      *uint                       `json:"parentCategoryId,omitempty"`
	SubCategories []*GetCategoriesTreeItemDto `json:"subCategories,omitempty"`
}

func MapGetCategoriesTreeItemsDto(categories []*entities.Category) []*GetCategoriesTreeItemDto {
	categoriesTreeItems := make([]*GetCategoriesTreeItemDto, len(categories))
	for categoryIndex, category := range categories {
		categoriesTreeItems[categoryIndex] = &GetCategoriesTreeItemDto{
			ID:            category.ID,
			Name:          category.Name,
			ParentID:      category.ParentCategoryID,
			SubCategories: MapGetCategoriesTreeItemsDto(category.Children),
		}
	}
	return categoriesTreeItems
}
