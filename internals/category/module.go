package category

import (
	"github.com/gin-gonic/gin"
	"github.com/ladmakhi81/learnup/internals/category/handler"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	categoryAdminHandler *handler.Handler
	middleware           *middleware.Middleware
	translationSvc       contracts.Translator
}

func NewModule(
	categoryAdminHandler *handler.Handler,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
) *Module {
	return &Module{
		categoryAdminHandler: categoryAdminHandler,
		middleware:           middleware,
		translationSvc:       translationSvc,
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	categoriesApi := api.Group("/categories")

	categoriesApi.Use(m.middleware.CheckAccessToken())

	categoriesApi.POST("/", utils.JsonHandler(m.translationSvc, m.categoryAdminHandler.CreateCategory))
	categoriesApi.GET("/tree", utils.JsonHandler(m.translationSvc, m.categoryAdminHandler.GetCategoriesTree))
	categoriesApi.GET("/page", utils.JsonHandler(m.translationSvc, m.categoryAdminHandler.GetCategories))
	categoriesApi.DELETE("/:categoryId", utils.JsonHandler(m.translationSvc, m.categoryAdminHandler.DeleteCategory))
}
