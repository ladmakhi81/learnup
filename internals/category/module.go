package category

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/category/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	categoryAdminHandler *handler.Handler
	middleware           *middleware.Middleware
}

func NewModule(
	categoryAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		categoryAdminHandler: categoryAdminHandler,
		middleware:           middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	categoriesApi := api.Group("/categories")

	categoriesApi.Use(m.middleware.CheckAccessToken())

	categoriesApi.POST("/", utils.JsonHandler(m.categoryAdminHandler.CreateCategory))
	categoriesApi.GET("/tree", utils.JsonHandler(m.categoryAdminHandler.GetCategoriesTree))
	categoriesApi.GET("/page", utils.JsonHandler(m.categoryAdminHandler.GetCategories))
	categoriesApi.DELETE("/:categoryId", utils.JsonHandler(m.categoryAdminHandler.DeleteCategory))
}
