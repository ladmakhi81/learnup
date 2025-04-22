package service

import (
	dtoreq "github.com/ladmakhi81/learnup/internals/course/dto/req"
	"github.com/ladmakhi81/learnup/internals/db"
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/types"
	"time"
)

type CourseService interface {
	Create(authContext any, dto dtoreq.CreateCourseReq) (*entities.Course, error)
	GetCourses(page, pageSize int) ([]*entities.Course, int, error)
	FindDetailById(id uint) (*entities.Course, error)
	VerifyCourse(authContext any, dto dtoreq.VerifyCourseReq) error
	UpdateIntroductionURL(dto dtoreq.UpdateIntroductionURLReq) error
	CreateCompleteIntroductionVideoNotification(id uint) error
}

type CourseServiceImpl struct {
	unitOfWork   db.UnitOfWork
	translateSvc contracts.Translator
}

func NewCourseServiceImpl(
	unitOfWork db.UnitOfWork,
	translateSvc contracts.Translator,
) *CourseServiceImpl {
	return &CourseServiceImpl{
		unitOfWork:   unitOfWork,
		translateSvc: translateSvc,
	}
}

func (svc CourseServiceImpl) Create(authContext any, dto dtoreq.CreateCourseReq) (*entities.Course, error) {
	const operationName = "CourseServiceImpl.Create"
	isCourseNameDuplicated, err := svc.unitOfWork.CourseRepo().Exist(
		map[string]any{"name": dto.Name},
	)
	if err != nil {
		return nil, types.NewServerError(
			"Error in checking course name exist or not",
			operationName,
			err,
		)
	}
	if isCourseNameDuplicated {
		return nil, types.NewConflictError(
			svc.translateSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, err := svc.unitOfWork.CategoryRepo().GetByID(dto.CategoryID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching category by id",
			operationName,
			err,
		)
	}
	if category == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_category")
	}
	teacher, err := svc.unitOfWork.UserRepo().GetByID(dto.TeacherID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching user teacher by id",
			operationName,
			err,
		)
	}
	if teacher == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_teacher")
	}
	authClaim := authContext.(*types.TokenClaim)
	authUser, err := svc.unitOfWork.UserRepo().GetByID(authClaim.UserID, nil)
	if err != nil {
		return nil, types.NewServerError(
			"Error in fetching logged in user",
			operationName,
			err,
		)
	}

	verifiedDate := time.Now()
	course := &entities.Course{
		CategoryID:                  &category.ID,
		Name:                        dto.Name,
		AbilityToAddComment:         dto.AbilityToAddComment,
		CanHaveDiscount:             dto.CanHaveDiscount,
		CommentAccessMode:           dto.CommentAccessMode,
		Description:                 dto.Description,
		DiscountFeeAmountPercentage: dto.DiscountFeeAmountPercentage,
		IsPublished:                 true,
		Image:                       dto.Image,
		IntroductionVideo:           dto.IntroductionVideo,
		IsVerifiedByAdmin:           true,
		Level:                       dto.Level,
		MaxDiscountAmount:           dto.MaxDiscountAmount,
		Prerequisite:                dto.Prerequisite,
		Status:                      entities.CourseStatus_Starting,
		Tags:                        dto.Tags,
		TeacherID:                   &teacher.ID,
		ThumbnailImage:              dto.ThumbnailImage,
		VerifiedDate:                &verifiedDate,
		VerifiedByID:                &authUser.ID,
	}
	course.SetPrice(dto.Price)
	course.SetFee(dto.Fee)

	if err := svc.unitOfWork.CourseRepo().Create(course); err != nil {
		return nil, types.NewServerError(
			"Create Course Throw Error",
			operationName,
			err,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) GetCourses(page, pageSize int) ([]*entities.Course, int, error) {
	const operationName = "CourseServiceImpl.GetCourses"
	courses, count, err := svc.unitOfWork.CourseRepo().GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Relations: []string{
			"Teacher",
			"Category",
			"VerifiedBy",
		},
	})
	if err != nil {
		return nil, 0, types.NewServerError(
			"Find All Pageable Courses Throw Error",
			operationName,
			err,
		)
	}
	return courses, count, nil
}

func (svc CourseServiceImpl) FindDetailById(id uint) (*entities.Course, error) {
	const operationName = "CourseServiceImpl.FindDetailById"
	course, err := svc.unitOfWork.CourseRepo().GetByID(id, []string{"Teacher", "Category", "VerifiedBy"})
	if err != nil {
		return nil, types.NewServerError(
			"Find Course Detail By ID Throw Error",
			operationName,
			err,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) VerifyCourse(authContext any, dto dtoreq.VerifyCourseReq) error {
	const operationName = "CourseServiceImpl.VerifyCourse"
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.ID, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
		)
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate(
				"course.errors.not_found",
			),
		)
	}
	if course.Status != entities.CourseStatus_InProgress {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.unable_to_verify"),
		)
	}
	adminClaim := authContext.(*types.TokenClaim)
	admin, err := svc.unitOfWork.UserRepo().GetByID(adminClaim.UserID, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching user admin by id",
			operationName,
			err,
		)
	}
	if admin == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("user.errors.admin_not_found"),
		)
	}
	if dto.Fee > course.Price || dto.Fee < 0 || dto.Fee > course.Price-course.MaxDiscountAmount {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.invalid_fee"),
		)
	}
	if dto.DiscountFeeAmountPercentage > 100 {
		return types.NewBadRequestError(
			svc.translateSvc.Translate("course.errors.invalid_max_discount_percentage"),
		)
	}
	now := time.Now()
	course.Fee = dto.Fee
	if course.CanHaveDiscount {
		course.DiscountFeeAmountPercentage = dto.DiscountFeeAmountPercentage
	}
	course.Status = entities.CourseStatus_Verified
	course.StatusChangedAt = &now
	course.VerifiedByID = &admin.ID
	course.VerifiedDate = &now
	course.IsVerifiedByAdmin = true
	if err := svc.unitOfWork.CourseRepo().Update(course); err != nil {
		return types.NewServerError(
			"Error in verifying the course by admin",
			operationName,
			err,
		)
	}
	//TODO: notification system
	// create notification for teacher that course verified
	// send email for this notification
	return nil
}

func (svc CourseServiceImpl) UpdateIntroductionURL(dto dtoreq.UpdateIntroductionURLReq) error {
	const operationName = "CourseServiceImpl.UpdateIntroductionURL"
	course, err := svc.unitOfWork.CourseRepo().GetByID(dto.CourseId, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
		)
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("course.errors.not_found"),
		)
	}
	course.IntroductionVideo = dto.URL
	if err := svc.unitOfWork.CourseRepo().Update(course); err != nil {
		return types.NewServerError(
			"Error in setting introduction video url",
			operationName,
			err,
		)
	}
	return nil
}

func (svc CourseServiceImpl) CreateCompleteIntroductionVideoNotification(id uint) error {
	const operationName = "CourseServiceImpl.CreateCompleteIntroductionVideoNotification"
	course, err := svc.unitOfWork.CourseRepo().GetByID(id, nil)
	if err != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			operationName,
			err,
		)
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("course.errors.not_found"),
		)
	}
	notification := &entities.Notification{
		Type: entities.NotificationType_CompleteIntroductionCourseVideoUpload,
		Metadata: map[string]any{
			"courseId":   course.ID,
			"courseName": course.Name,
		},
		IsSeen: false,
		UserID: course.TeacherID,
	}
	if err := svc.unitOfWork.NotificationRepo().Create(notification); err != nil {
		return types.NewServerError(
			"Error in creating notification",
			operationName,
			err,
		)
	}
	return nil
}
