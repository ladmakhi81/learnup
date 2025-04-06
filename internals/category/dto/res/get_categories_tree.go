package dtores

import "github.com/ladmakhi81/learnup/internals/category/entity"

type GetCategoriesTreeItem struct {
	ID            uint                     `json:"id"`
	Name          string                   `json:"name"`
	ParentID      *uint                    `json:"parentCategoryId,omitempty"`
	SubCategories []*GetCategoriesTreeItem `json:"subCategories,omitempty"`
}

type GetCategoriesTreeRes struct {
	Categories []*GetCategoriesTreeItem `json:"categories"`
}

func mapCategoryToCategoryTreeItem(category []*entity.Category) []*GetCategoriesTreeItem {
	categoriesTreeItems := make([]*GetCategoriesTreeItem, len(category))
	for categoryIndex, category := range category {
		categoriesTreeItems[categoryIndex] = &GetCategoriesTreeItem{
			ID:            category.ID,
			Name:          category.Name,
			ParentID:      category.ParentCategoryID,
			SubCategories: mapCategoryToCategoryTreeItem(category.Children),
		}
	}
	return categoriesTreeItems
}

func NewGetCategoriesTreeRes(categories []*entity.Category) GetCategoriesTreeRes {
	categoriesTree := GetCategoriesTreeRes{}
	categoriesTree.Categories = mapCategoryToCategoryTreeItem(categories)
	return categoriesTree
}
