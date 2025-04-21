package handler

import (
	"github.com/gin-gonic/gin"
	dtores "github.com/ladmakhi81/learnup/internals/teacher/dto/res"
	teacherService "github.com/ladmakhi81/learnup/internals/teacher/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type QuestionHandler struct {
	translationSvc     contracts.Translator
	teacherQuestionSvc teacherService.TeacherQuestionService
}

func NewQuestionHandler(
	translationSvc contracts.Translator,
	teacherQuestionSvc teacherService.TeacherQuestionService,
) *QuestionHandler {
	return &QuestionHandler{
		translationSvc:     translationSvc,
		teacherQuestionSvc: teacherQuestionSvc,
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
	authContext, _ := ctx.Get("AUTH")
	teacherClaim := authContext.(*types.TokenClaim)
	page, pageSize := utils.ExtractPaginationMetadata(ctx.Query("page"), ctx.Query("pageSize"))
	questions, questionsErr := h.teacherQuestionSvc.GetQuestions(teacherService.GetQuestionOptions{
		PageSize:  pageSize,
		Page:      page,
		CourseID:  courseID,
		TeacherID: teacherClaim.UserID,
	})
	if questionsErr != nil {
		return nil, questionsErr
	}
	questionsCount, questionsCountErr := h.teacherQuestionSvc.GetQuestionCount(courseID)
	if questionsCountErr != nil {
		return nil, questionsCountErr
	}
	questionsRes := types.NewPaginationRes(
		dtores.MapGetQuestionItemRes(questions),
		page,
		utils.CalculatePaginationTotalPage(questionsCount, pageSize),
		questionsCount,
	)
	return types.NewApiResponse(http.StatusOK, questionsRes), nil
}
