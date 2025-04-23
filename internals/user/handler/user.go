package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/dto/res"
	"github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"net/http"
)

type Handler struct {
	userSvc        service.UserSvc
	validationSvc  contracts.Validation
	translationSvc contracts.Translator
}

func NewHandler(
	userSvc service.UserSvc,
	validationSvc contracts.Validation,
	translationSvc contracts.Translator,
) *Handler {
	return &Handler{
		userSvc:        userSvc,
		validationSvc:  validationSvc,
		translationSvc: translationSvc,
	}
}

// CreateBasicUser godoc
//
//	@Summary	Create Basic User
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		CreateBasicUserReqDto	body		dtoreq.CreateBasicUserReqDto	true	" "
//	@Success	201					{object}	types.ApiResponse{data=dtores.CreateBasicUserResDto}
//	@Failure	400					{object}	types.ApiError
//	@Failure	409					{object}	types.ApiError
//	@Failure	500					{object}	types.ApiError
//	@Router		/users/basic [post]
//
// @Security BearerAuth
func (h Handler) CreateBasicUser(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.CreateBasicUserReqDto)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, err := h.userSvc.CreateBasic(*dto)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusCreated, dtores.NewCreateBasicUserResDto(user)), nil
}
