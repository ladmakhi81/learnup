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
	repo         *db.Repositories
	translateSvc contracts.Translator
}

func NewCourseServiceImpl(
	repo *db.Repositories,
	translateSvc contracts.Translator,
) *CourseServiceImpl {
	return &CourseServiceImpl{
		repo:         repo,
		translateSvc: translateSvc,
	}
}

func (svc CourseServiceImpl) Create(authContext any, dto dtoreq.CreateCourseReq) (*entities.Course, error) {
	isCourseNameDuplicated, courseNameDuplicatedErr := svc.repo.CourseRepo.Exist(
		map[string]any{"name": dto.Name},
	)
	if courseNameDuplicatedErr != nil {
		return nil, types.NewServerError(
			"Error in checking course name exist or not",
			"CourseServiceImpl.Create",
			courseNameDuplicatedErr,
		)
	}
	if isCourseNameDuplicated {
		return nil, types.NewConflictError(
			svc.translateSvc.Translate("course.errors.name_duplicate"),
		)
	}
	category, categoryErr := svc.repo.CategoryRepo.GetByID(dto.CategoryID, nil)
	if categoryErr != nil {
		return nil, types.NewServerError(
			"Error in fetching category by id",
			"CourseServiceImpl.Create",
			categoryErr,
		)
	}
	if category == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_category")
	}
	teacher, teacherErr := svc.repo.UserRepo.GetByID(dto.TeacherID, nil)
	if teacherErr != nil {
		return nil, types.NewServerError(
			"Error in fetching user teacher by id",
			"CourseServiceImpl.Create",
			teacherErr,
		)
	}
	if teacher == nil {
		return nil, types.NewNotFoundError("course.errors.not_found_teacher")
	}
	authClaim := authContext.(*types.TokenClaim)
	authUser, authUserErr := svc.repo.UserRepo.GetByID(authClaim.UserID, nil)
	if authUserErr != nil {
		return nil, types.NewServerError(
			"Error in fetching logged in user",
			"CourseServiceImpl.Create",
			authUserErr,
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

	if err := svc.repo.CourseRepo.Create(course); err != nil {
		return nil, types.NewServerError(
			"Create Course Throw Error",
			"CourseServiceImpl.Create",
			err,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) GetCourses(page, pageSize int) ([]*entities.Course, int, error) {
	courses, count, coursesErr := svc.repo.CourseRepo.GetPaginated(repositories.GetPaginatedOptions{
		Offset: &page,
		Limit:  &pageSize,
		Relations: []string{
			"Teacher",
			"Category",
			"VerifiedBy",
		},
	})
	if coursesErr != nil {
		return nil, 0, types.NewServerError(
			"Find All Pageable Courses Throw Error",
			"CourseServiceImpl.FetchPage",
			coursesErr,
		)
	}
	return courses, count, nil
}

func (svc CourseServiceImpl) FindDetailById(id uint) (*entities.Course, error) {
	course, courseErr := svc.repo.CourseRepo.GetByID(id, []string{"Teacher", "Category", "VerifiedBy"})
	if courseErr != nil {
		return nil, types.NewServerError(
			"Find Course Detail By ID Throw Error",
			"CourseServiceImpl.FindDetailByID",
			courseErr,
		)
	}
	return course, nil
}

func (svc CourseServiceImpl) VerifyCourse(authContext any, dto dtoreq.VerifyCourseReq) error {
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.ID, nil)
	if courseErr != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			"CourseServiceImpl.VerifyCourse",
			courseErr,
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
	admin, adminErr := svc.repo.UserRepo.GetByID(adminClaim.UserID, nil)
	if adminErr != nil {
		return types.NewServerError(
			"Error in fetching user admin by id",
			"CourseServiceImpl.VerifyCourse",
			adminErr,
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
	if err := svc.repo.CourseRepo.Update(course); err != nil {
		return types.NewServerError(
			"Error in verifying the course by admin",
			"CourseServiceImpl.VerifyCourse",
			err,
		)
	}
	//TODO: notification system
	// create notification for teacher that course verified
	// send email for this notification
	return nil
}

func (svc CourseServiceImpl) UpdateIntroductionURL(dto dtoreq.UpdateIntroductionURLReq) error {
	course, courseErr := svc.repo.CourseRepo.GetByID(dto.CourseId, nil)
	if courseErr != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			"CourseServiceImpl.UpdateIntroductionURL",
			courseErr,
		)
	}
	if course == nil {
		return types.NewNotFoundError(
			svc.translateSvc.Translate("course.errors.not_found"),
		)
	}
	course.IntroductionVideo = dto.URL
	if err := svc.repo.CourseRepo.Update(course); err != nil {
		return types.NewServerError(
			"Error in setting introduction video url",
			"CourseServiceImpl.UpdateIntroductionURL",
			err,
		)
	}
	return nil
}

func (svc CourseServiceImpl) CreateCompleteIntroductionVideoNotification(id uint) error {
	course, courseErr := svc.repo.CourseRepo.GetByID(id, nil)
	if courseErr != nil {
		return types.NewServerError(
			"Error in fetching course by id",
			"CourseServiceImpl.CreateCompleteIntroductionVideoNotification",
			courseErr,
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
	notificationErr := svc.repo.NotificationRepo.Create(notification)
	if notificationErr != nil {
		return types.NewServerError(
			"Error in creating notification",
			"CourseServiceImpl.CreateCompleteIntroductionVideoNotification",
			notificationErr,
		)
	}
	return nil
}
