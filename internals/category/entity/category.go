package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name             string      `gorm:"type:varchar(250);not null;column:name;unique"`
	IsPublished      bool        `gorm:"type:boolean;default:false;not null;column:is_published"`
	ParentCategoryID *uint       `gorm:"index;column:parent_category_id;type:int unsigned;"`
	ParentCategory   *Category   `gorm:"foreignKey:ParentCategoryID"`
	Children         []*Category `gorm:"foreignKey:ParentCategoryID"`
}

func (Category) TableName() string {
	return "_categories"
}
