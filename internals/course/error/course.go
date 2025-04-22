package error

import "github.com/ladmakhi81/learnup/types"

var (
	Course_NotFound                     = types.NewNotFoundError("course.errors.not_found")
	Course_NameDuplicated               = types.NewConflictError("course.errors.name_duplicate")
	Course_NotFoundCategory             = types.NewNotFoundError("course.errors.not_found_category")
	Course_NotFoundTeacher              = types.NewNotFoundError("course.errors.not_found_teacher")
	Course_UnableToVerify               = types.NewBadRequestError("course.errors.unable_to_verify")
	Course_InvalidFee                   = types.NewBadRequestError("course.errors.invalid_fee")
	Course_InvalidMaxDiscountPercentage = types.NewBadRequestError("course.errors.invalid_max_discount_percentage")
	Course_ForbiddenAccess              = types.NewForbiddenAccessError("common.errors.forbidden_access")
)
