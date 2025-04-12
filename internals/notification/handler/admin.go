package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/notification/dto/res"
	notificationService "github.com/ladmakhi81/learnup/internals/notification/service"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
	"strconv"
)

type NotificationAdminHandler struct {
	notificationSvc notificationService.NotificationService
}

func NewNotificationAdminHandler(notificationSvc notificationService.NotificationService) *NotificationAdminHandler {
	return &NotificationAdminHandler{
		notificationSvc: notificationSvc,
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
//	@Router		/notifications/admin/{notification-id}/seen [patch]
//
//	@Security	BearerAuth
func (h NotificationAdminHandler) SeenNotification(ctx *gin.Context) (*types.ApiResponse, error) {
	notificationIdParam := ctx.Param("notification-id")
	notificationID, notificationIDErr := strconv.Atoi(notificationIdParam)
	if notificationIDErr != nil {
		return nil, types.NewBadRequestError("invalid notification id")
	}
	if err := h.notificationSvc.SeenById(uint(notificationID)); err != nil {
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
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.NotificationPageItem}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/notifications/admin/page [get]
//
//	@Security	BearerAuth
func (h NotificationAdminHandler) GetNotificationsPage(ctx *gin.Context) (*types.ApiResponse, error) {
	page, pageSize := utils.ExtractPaginationMetadata(
		ctx.Query("page"),
		ctx.Query("pageSize"),
	)
	notifications, notificationsErr := h.notificationSvc.FetchPageable(page, pageSize)
	if notificationsErr != nil {
		return nil, notificationsErr
	}
	notificationsCount, notificationsCountErr := h.notificationSvc.FetchCount()
	if notificationsCountErr != nil {
		return nil, notificationsCountErr
	}
	mappedNotifications := dtores.NewNotificationPageItems(notifications)
	res := types.NewPaginationRes(
		mappedNotifications,
		page,
		notificationsCount,
		utils.CalculatePaginationTotalPage(notificationsCount),
	)
	return types.NewApiResponse(http.StatusOK, res), nil
}
