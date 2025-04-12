package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/auth/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/auth/dto/res"
	"github.com/ladmakhi81/learnup/internals/auth/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
)

type Handler struct {
	authSvc        service.AuthService
	validationSvc  validation.Validation
	translationSvc translations.Translator
}

func NewHandler(
	authSvc service.AuthService,
	validationSvc validation.Validation,
	translationSvc translations.Translator,
) *Handler {
	return &Handler{
		authSvc:        authSvc,
		validationSvc:  validationSvc,
		translationSvc: translationSvc,
	}
}

// Login godoc
//
//	@Summary	Login a user and return an access token
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		loginRequest	body		dtoreq.LoginReq	true	" "
//	@Success	200				{object}	types.ApiResponse{data=dtores.LoginRes}
//	@Failure	400				{object}	types.ApiError
//	@Failure	404				{object}	types.ApiError
//	@Failure	500				{object}	types.ApiError
//	@Router		/auth/login [post]
func (h Handler) Login(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := new(dtoreq.LoginReq)
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	accessToken, accessTokenErr := h.authSvc.Login(*dto)
	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	loginRes := dtores.NewLoginRes(accessToken)
	return types.NewApiResponse(http.StatusOK, loginRes), nil
}
