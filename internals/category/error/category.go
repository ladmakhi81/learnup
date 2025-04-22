package error

import "github.com/ladmakhi81/learnup/types"

var (
	Category_DuplicateName  = types.NewConflictError("category.errors.name_duplicate")
	Category_ParentNotFound = types.NewNotFoundError("category.errors.parent_category_id_not_found")
	Category_NotFound       = types.NewNotFoundError("category.errors.not_found")
)
