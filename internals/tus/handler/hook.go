package handler

import (
	"github.com/gin-gonic/gin"
	reqdto "github.com/ladmakhi81/learnup/internals/tus/dto"
	"github.com/ladmakhi81/learnup/internals/tus/service"
	"github.com/ladmakhi81/learnup/types"
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
	dto := &reqdto.TusWebhookDTO{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewServerError(
			"Error in convert data into tus dto",
			"TusHookHandler.VideoWebhook",
			err,
		)
	}
	h.tusHookSvc.VideoWebhook(*dto)
	return types.NewApiResponse(http.StatusOK, gin.H{"message": "hook received"}), nil
}
