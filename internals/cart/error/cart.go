package error

import (
	"github.com/ladmakhi81/learnup/shared/types"
)

var (
	Cart_Duplicated      = types.NewConflictError("cart.errors.exist_before")
	Cart_NotFound        = types.NewNotFoundError("cart.errors.not_found")
	Cart_ForbiddenAccess = types.NewForbiddenAccessError("cart.errors.owner_delete")
	Cart_ListNotMatch    = types.NewNotFoundError("cart.errors.list_not_match")
)
