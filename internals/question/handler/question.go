package handler

import (
	"github.com/gin-gonic/gin"
	questionDtoReq "github.com/ladmakhi81/learnup/internals/question/dto/req"
	questionDtoRes "github.com/ladmakhi81/learnup/internals/question/dto/res"
	answerService "github.com/ladmakhi81/learnup/internals/question/service"
	userService "github.com/ladmakhi81/learnup/internals/user/service"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"github.com/ladmakhi81/learnup/utils"
	"net/http"
)

type Handler struct {
	answerSvc      answerService.QuestionAnswerService
	translationSvc contracts.Translator
	validationSvc  contracts.Validation
	userSvc        userService.UserSvc
}

func NewHandler(
	answerSvc answerService.QuestionAnswerService,
	translationSvc contracts.Translator,
	validationSvc contracts.Validation,
	userSvc userService.UserSvc,
) *Handler {
	return &Handler{
		answerSvc:      answerSvc,
		translationSvc: translationSvc,
		validationSvc:  validationSvc,
		userSvc:        userSvc,
	}
}

// AnswerQuestion godoc
//
//	@Summary	Submit an answer to a question
//	@Tags		questions
//	@Accept		json
//	@Produce	json
//	@Param		question-id	path		uint								true	" "
//	@Param		answer		body		questionDtoReq.AnswerQuestionReqDto	true	" "
//	@Success	201			{object}	types.ApiResponse{data=questionDtoRes.CreateAnswerResDto}
//	@Failure	400			{object}	types.ApiResponse
//	@Failure	401			{object}	types.ApiResponse
//	@Failure	404			{object}	types.ApiResponse
//	@Failure	500			{object}	types.ApiResponse
//	@Router		/questions/{question-id}/answer [post]
//	@Security	BearerAuth
func (h Handler) AnswerQuestion(ctx *gin.Context) (*types.ApiResponse, error) {
	user, err := h.userSvc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	questionID, err := utils.ToUint(ctx.Param("question-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("question.errors.invalid_id"),
		)
	}
	dto := &questionDtoReq.AnswerQuestionReqDto{}
	if err := ctx.Bind(dto); err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("common.errors.invalid_request_body"),
		)
	}
	if err := h.validationSvc.Validate(dto); err != nil {
		return nil, err
	}
	dto.QuestionID = questionID
	answer, err := h.answerSvc.Create(user, *dto)
	if err != nil {
		return nil, err
	}
	answerRes := questionDtoRes.NewCreateAnswerResDto(answer)
	return types.NewApiResponse(http.StatusCreated, answerRes), nil
}

// GetQuestionAnswers godoc
//
//	@Summary	Retrieve answers for a specific question
//	@Tags		questions
//	@Accept		json
//	@Produce	json
//	@Param		question-id	path		uint	true	" "
//	@Success	200			{object}	types.ApiResponse{data=[]questionDtoRes.GetAnswerItemDto}
//	@Failure	400			{object}	types.ApiError
//	@Failure	401			{object}	types.ApiError
//	@Failure	404			{object}	types.ApiError
//	@Failure	500			{object}	types.ApiError
//	@Router		/questions/{question-id}/answers [get]
//	@Security	BearerAuth
func (h Handler) GetQuestionAnswers(ctx *gin.Context) (*types.ApiResponse, error) {
	questionID, err := utils.ToUint(ctx.Param("question-id"))
	if err != nil {
		return nil, types.NewBadRequestError(
			h.translationSvc.Translate("question.errors.invalid_id"),
		)
	}
	answers, err := h.answerSvc.GetQuestionAnswers(questionID)
	if err != nil {
		return nil, err
	}
	return types.NewApiResponse(http.StatusOK, questionDtoRes.MapGetAnswerItemsDto(answers)), nil
}
