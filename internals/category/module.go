package category

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/category/handler"
	"github.com/ladmakhi81/learnup/internals/middleware"
	"github.com/ladmakhi81/learnup/utils"
)

type Module struct {
	categoryAdminHandler *handler.CategoryAdminHandler
	middleware           *middleware.Middleware
}

func NewModule(
	categoryAdminHandler *handler.CategoryAdminHandler,
	middleware *middleware.Middleware,
) *Module {
	return &Module{
		categoryAdminHandler: categoryAdminHandler,
		middleware:           middleware,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	categoriesApi := api.Group("/categories")
	categoriesAdminApi := categoriesApi.Group("/admin")

	categoriesAdminApi.Use(m.middleware.CheckAccessToken())

	categoriesAdminApi.POST("/", utils.JsonHandler(m.categoryAdminHandler.CreateCategory))
	categoriesAdminApi.GET("/tree", utils.JsonHandler(m.categoryAdminHandler.GetCategoriesTree))
	categoriesAdminApi.GET("/page", utils.JsonHandler(m.categoryAdminHandler.GetCategories))
	categoriesAdminApi.DELETE("/:categoryId", utils.JsonHandler(m.categoryAdminHandler.DeleteCategory))
}
