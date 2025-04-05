package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/user/constant"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/dto/res"
	"github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
)

type UserAdminHandler struct {
	userSvc       service.UserSvc
	validationSvc validation.Validation
}

func NewUserAdminHandler(
	userSvc service.UserSvc,
	validationSvc validation.Validation,
) *UserAdminHandler {
	return &UserAdminHandler{
		userSvc:       userSvc,
		validationSvc: validationSvc,
	}
}

// CreateBasicUser godoc
//
//	@Summary	Create Basic User
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		CreateBasicUserReq	body		dtoreq.CreateBasicUserReq	true	" "
//	@Success	201					{object}	res.CreateBasicUserRes
//	@Failure	400					{object}	types.ApiError
//	@Failure	409					{object}	types.ApiError
//	@Failure	500					{object}	types.ApiError
//	@Router		/users/admin/basic [post]
func (h UserAdminHandler) CreateBasicUser(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.CreateBasicUserReq)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(constant.InvalidRequestBody)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, userErr := h.userSvc.CreateBasic(*dto)
	if userErr != nil {
		return nil, userErr
	}
	userResponse := res.NewCreateUserResponse(user)
	return types.NewApiResponse(http.StatusCreated, userResponse), nil
}
