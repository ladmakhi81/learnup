package handler

import (
	"github.com/gin-gonic/gin"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type QuestionHandler struct {
	translationSvc     contracts.Translator
	teacherQuestionSvc teacherService.TeacherQuestionService
	userSvc            userService.UserSvc
}

func NewQuestionHandler(
	translationSvc contracts.Translator,
	teacherQuestionSvc teacherService.TeacherQuestionService,
	userSvc userService.UserSvc,
) *QuestionHandler {
	return &QuestionHandler{
		translationSvc:     translationSvc,
		teacherQuestionSvc: teacherQuestionSvc,
		userSvc:            userSvc,
	}
}

// GetQuestions godoc
//
//	@Summary	Get questions by course ID for teacher
//	@Tags		teacher
//	@Accept		json
//	@Produce	json
//	@Param		course-id	query		uint	false	"Course ID"
//	@Param		page		query		int		false	"Page number"		default(0)
//	@Param		pageSize	query		int		false	"Items per page"	default(10)
//	@Success	200			{object}	types.ApiResponse{data=types.PaginationRes{row=[]dtores.GetQuestionItemRes}}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Security	BearerAuth
//	@Router		/teacher/questions [get]
func (h QuestionHandler) GetQuestions(ctx *gin.Context) (*types.ApiResponse, error) {
	courseIDParam := ctx.Query("course-id")
	var courseID *uint
	if courseIDParam != "" {
		parsedCourseID, parsedCourseIDErr := utils.ToUint(courseIDParam)
		if parsedCourseIDErr != nil {
			return nil, types.NewBadRequestError(
				h.translationSvc.Translate("course.errors.invalid_course_id"),
			)
		}
		courseID = &parsedCourseID
	}
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	questions, count, questionsErr := h.teacherQuestionSvc.GetQuestions(user, teacherService.GetQuestionOptions{PageSize: pageSize, Page: page, CourseID: courseID})
	if questionsErr != nil {
		return nil, questionsErr
	}
	questionsRes := types.NewPaginationRes(
		dtores.MapGetQuestionItemRes(questions),
		page,
		utils.CalculatePaginationTotalPage(count, pageSize),
		count,
	)
	return types.NewApiResponse(http.StatusOK, questionsRes), nil
}
