package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/notification/dto/res"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/types"
	"github.com/ladmakhi81/learnup/shared/utils"
	"net/http"
)

type Handler struct {
	notificationSvc notificationService.NotificationService
	translationSvc  contracts.Translator
}

func NewHandler(
	notificationSvc notificationService.NotificationService,
	translationSvc contracts.Translator,
) *Handler {
	return &Handler{
		notificationSvc: notificationSvc,
		translationSvc:  translationSvc,
	}
}

// SeenNotification godoc
//
//	@Summary	Mark notification as seen
//	@Tags		notifications
//	@Param		notification-id	path		int	true	"Notification ID"
//	@Success	200				{object}	types.ApiResponse
//	@Failure	400				{object}	types.ApiError
//	@Failure	404				{object}	types.ApiError
//	@Failure	500				{object}	types.ApiError
//	@Router		/notifications/{notification-id}/seen [patch]
//
//	@Security	BearerAuth
func (h Handler) SeenNotification(ctx *gin.Context) (*types.ApiResponse, error) {
	notificationID, err := utils.ToUint(ctx.Param("notification-id"))
	if err != nil {
		return nil, types.NewBadRequestError(h.translationSvc.Translate("notifications.errors.invalid_notification_id"))
	}
	if err := h.notificationSvc.SeenById(notificationID); err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, map[string]any{}), nil
}

// GetNotificationsPage godoc
//
//	@Summary	Get paginated notifications
//	@Tags		notifications
//	@Param		page		query		int	false	"Page number"		default(0)
//	@Param		pageSize	query		int	false	"Number per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.NotificationPageItemDto}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/notifications/page [get]
//
//	@Security	BearerAuth
func (h Handler) GetNotificationsPage(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	notifications, count, err := h.notificationSvc.FetchPageable(page, pageSize)
	if err != nil {
		return nil, err
	}
	res := types.NewPaginationRes(
		dtores.NewNotificationPageItemsDto(notifications),
		page,
		count,
		utils.CalculatePaginationTotalPage(count, pageSize),
	)
	return types.NewApiResponse(http.StatusOK, res), nil
}
