package handler

import (
	"github.com/gin-gonic/gin"
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	"github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/shared/types"
	"net/http"
)

type TusHookHandler struct {
	tusHookSvc service.TusService
}

func NewTusHookHandler(tusHookSvc service.TusService) *TusHookHandler {
	return &TusHookHandler{
		tusHookSvc: tusHookSvc,
	}
}

func (h TusHookHandler) VideoWebhook(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &reqdto.TusWebhookDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewServerError(
			"Error in convert data into tus dto",
			"TusHookHandler.VideoWebhook",
			err,
		)
	}
	switch dto.Type {
	case reqdto.TusHookType_PostFinish:
		go h.tusHookSvc.VideoWebhook(ctx, *dto)
	}
	return types.NewApiResponse(http.StatusOK, gin.H{"message": "hook received"}), nil
}
