package category

import (
	"github.com/gin-gonic/gin"
	categoryHandler "github.com/ladmakhi81/learnup/internals/category/handler"
	categoryService "github.com/ladmakhi81/learnup/internals/category/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/shared/middleware"
	"github.com/ladmakhi81/learnup/shared/utils"
)

type Module struct {
	categoryHandler *categoryHandler.Handler
	middleware      *middleware.Middleware
	translationSvc  contracts.Translator
}

func NewModule(
	categorySvc categoryService.CategoryService,
	middleware *middleware.Middleware,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
) *Module {
	return &Module{
		middleware:     middleware,
		translationSvc: translationSvc,
		categoryHandler: categoryHandler.NewHandler(
			categorySvc,
			translationSvc,
			validationSvc,
		),
	}
}

func (m Module) Register(api *gin.RouterGroup) {
	categoriesApi := api.Group("/categories")
	categoriesApi.Use(m.middleware.CheckAccessToken())
	categoriesApi.POST("/", utils.JsonHandler(m.translationSvc, m.categoryHandler.CreateCategory))
	categoriesApi.GET("/tree", utils.JsonHandler(m.translationSvc, m.categoryHandler.GetCategoriesTree))
	categoriesApi.GET("/page", utils.JsonHandler(m.translationSvc, m.categoryHandler.GetCategories))
	categoriesApi.DELETE("/:categoryId", utils.JsonHandler(m.translationSvc, m.categoryHandler.DeleteCategory))
}
