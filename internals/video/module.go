package video

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/video/handler"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	videoAdminHandler *handler.VideoAdminHandler
}

func NewModule(
	videoAdminHandler *handler.VideoAdminHandler,
) *Module {
	return &Module{
		videoAdminHandler: videoAdminHandler,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	videosApi := api.Group("/videos")
	videosAdminApi := videosApi.Group("/admin")

	videosAdminApi.POST("/", utils.JsonHandler(m.videoAdminHandler.AddNewVideoToCourse))
}
