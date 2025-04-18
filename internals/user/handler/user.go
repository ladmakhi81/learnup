package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/user/dto/req"
	"github.com/ladmakhi81/learnup/internals/user/dto/res"
	"github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
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
//	@Param		CreateBasicUserReq	body		dtoreq.CreateBasicUserReq	true	" "
//	@Success	201					{object}	types.ApiResponse{data=dtores.CreateBasicUserRes}
//	@Failure	400					{object}	types.ApiError
//	@Failure	409					{object}	types.ApiError
//	@Failure	500					{object}	types.ApiError
//	@Router		/users/basic [post]
//
// @Security BearerAuth
func (h Handler) CreateBasicUser(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.CreateBasicUserReq)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	user, userErr := h.userSvc.CreateBasic(*dto)
	if userErr != nil {
		return nil, userErr
	}
	userResponse := dtores.NewCreateUserResponse(user)
	return types.NewApiResponse(http.StatusCreated, userResponse), nil
}
