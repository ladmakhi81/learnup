package dtores

import (
	"github.com/ladmakhi81/learnup/db/entities"
)

type GetCategoriesTreeItem struct {
	ID            uint                     `json:"id"`
	Name          string                   `json:"name"`
	ParentID      *uint                    `json:"parentCategoryId,omitempty"`
	SubCategories []*GetCategoriesTreeItem `json:"subCategories,omitempty"`
}

type GetCategoriesTreeRes struct {
	Categories []*GetCategoriesTreeItem `json:"categories"`
}

func mapCategoryToCategoryTreeItem(categories []*entities.Category) []*GetCategoriesTreeItem {
	categoriesTreeItems := make([]*GetCategoriesTreeItem, len(categories))
	for categoryIndex, category := range categories {
		categoriesTreeItems[categoryIndex] = &GetCategoriesTreeItem{
			ID:            category.ID,
			Name:          category.Name,
			ParentID:      category.ParentCategoryID,
			SubCategories: mapCategoryToCategoryTreeItem(category.Children),
		}
	}
	return categoriesTreeItems
}

func NewGetCategoriesTreeRes(categories []*entities.Category) GetCategoriesTreeRes {
	categoriesTree := GetCategoriesTreeRes{}
	categoriesTree.Categories = mapCategoryToCategoryTreeItem(categories)
	return categoriesTree
}
