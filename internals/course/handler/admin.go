package handler

import (
	"github.com/gin-gonic/gin"
	dtoreq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	dtores "github.com/ladmakhi81/learnup/internals/course/dto/res"
	"github.com/ladmakhi81/learnup/internals/course/service"
	"github.com/ladmakhi81/learnup/pkg/translations"
	"github.com/ladmakhi81/learnup/pkg/validation"
	"github.com/ladmakhi81/learnup/types"
	"net/http"
)

type CourseAdminHandler struct {
	courseSvc     service.CourseService
	validationSvc validation.Validation
	translateSvc  translations.Translator
}

func NewCourseAdminHandler(
	courseSvc service.CourseService,
	validationSvc validation.Validation,
	translateSvc translations.Translator,
) *CourseAdminHandler {
	return &CourseAdminHandler{
		courseSvc:     courseSvc,
		validationSvc: validationSvc,
		translateSvc:  translateSvc,
	}
}

// CreateCourse godoc
//
//	@Summary	Create a new course
//	@Tags		courses
//	@Accept		json
//	@Produce	json
//	@Param		requestBody	body		dtoreq.CreateCourseReq	true	" "
//	@Success	201			{object}	types.ApiResponse{data=dtores.CreateCourseRes}
//	@Failure	400			{object}	types.ApiResponse
//	@Failure	401			{object}	types.ApiResponse
//	@Failure	409			{object}	types.ApiResponse
//	@Failure	500			{object}	types.ApiResponse
//	@Router		/courses/admin [post]
//	@Security	BearerAuth
func (h CourseAdminHandler) CreateCourse(ctx *gin.Context) (*types.ApiResponse, error) {
	dto := &dtoreq.CreateCourseReq{}
	if err := ctx.ShouldBind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translateSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	authContext, _ := ctx.Get("AUTH")
	course, courseErr := h.courseSvc.Create(authContext, *dto)
	if courseErr != nil {
		return nil, courseErr
	}
	courseRes := dtores.CreateCourseRes{
		ID:                          course.ID,
		Fee:                         course.Fee,
		Price:                       course.Price,
		VerifiedByID:                course.VerifiedByID,
		VerifiedDate:                course.VerifiedDate,
		TeacherID:                   course.TeacherID,
		ThumbnailImage:              course.ThumbnailImage,
		Tags:                        course.Tags,
		Status:                      course.Status,
		Prerequisite:                course.Prerequisite,
		MaxDiscountAmount:           course.MaxDiscountAmount,
		Level:                       course.Level,
		IsVerifiedByAdmin:           course.IsVerifiedByAdmin,
		IntroductionVideo:           course.IntroductionVideo,
		Image:                       course.Image,
		IsPublished:                 course.IsPublished,
		DiscountFeeAmountPercentage: course.DiscountFeeAmountPercentage,
		Description:                 course.Description,
		CommentAccessMode:           course.CommentAccessMode,
		CanHaveDiscount:             course.CanHaveDiscount,
		AbilityToAddComment:         course.AbilityToAddComment,
		Name:                        course.Name,
		CategoryID:                  course.CategoryID,
	}
	return types.NewApiResponse(http.StatusCreated, courseRes), nil
}
