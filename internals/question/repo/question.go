package repo

import (
	"errors"
	"github.com/ladmakhi81/learnup/db"
	"github.com/ladmakhi81/learnup/db/entities"
	"gorm.io/gorm"
)

type FetchQuestionOptions struct {
	Page     *int
	PageSize *int
	UserID   *uint
	CourseID *uint
	Priority *entities.QuestionPriority
	IsClosed *bool
	VideoID  *uint
}

type FetchCountQuestionOptions struct {
	UserID   *uint
	CourseID *uint
	Priority *entities.QuestionPriority
	IsClosed *bool
	VideoID  *uint
}

type FetchOneQuestionOptions struct {
	UserID   *uint
	CourseID *uint
	VideoID  *uint
	ID       *uint
}

type QuestionRepo interface {
	Create(question *entities.Question) error
	Update(question *entities.Question) error
	Fetch(options FetchQuestionOptions) ([]*entities.Question, error)
	FetchCount(options FetchCountQuestionOptions) (int, error)
	Delete(id uint) error
	FindOne(options FetchOneQuestionOptions) (*entities.Question, error)
}

type QuestionRepoImpl struct {
	dbClient *db.Database
}

func NewQuestionRepoImpl(dbClient *db.Database) *QuestionRepoImpl {
	return &QuestionRepoImpl{
		dbClient: dbClient,
	}
}

func (repo *QuestionRepoImpl) Create(question *entities.Question) error {
	tx := repo.dbClient.Core.Create(question)
	return tx.Error
}

func (repo *QuestionRepoImpl) Update(question *entities.Question) error {
	tx := repo.dbClient.Core.Updates(question)
	return tx.Error
}

func (repo *QuestionRepoImpl) Fetch(options FetchQuestionOptions) ([]*entities.Question, error) {
	var questions []*entities.Question
	query := repo.dbClient.Core
	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	if options.Priority != nil {
		query = query.Where("priority = ?", options.Priority)
	}
	if options.IsClosed != nil {
		query = query.Where("is_closed = ?", options.IsClosed)
	}
	if options.VideoID != nil {
		query = query.Where("video_id = ?", options.VideoID)
	}
	if options.PageSize != nil && options.Page != nil {
		query = query.Offset((*options.Page) * (*options.PageSize)).Limit(*options.PageSize)
	}
	tx := query.Preload("User").Preload("Course").Preload("Video").Order("created_at desc").Find(&questions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return questions, nil
}

func (repo *QuestionRepoImpl) FetchCount(options FetchCountQuestionOptions) (int, error) {
	var count int64
	query := repo.dbClient.Core.Model(&entities.Question{})
	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	if options.Priority != nil {
		query = query.Where("priority = ?", options.Priority)
	}
	if options.IsClosed != nil {
		query = query.Where("is_closed = ?", options.IsClosed)
	}
	if options.VideoID != nil {
		query = query.Where("video_id = ?", options.VideoID)
	}
	tx := query.Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo *QuestionRepoImpl) Delete(id uint) error {
	tx := repo.dbClient.Core.Delete(id)
	return tx.Error
}

func (repo *QuestionRepoImpl) FindOne(options FetchOneQuestionOptions) (*entities.Question, error) {
	question := &entities.Question{}
	query := repo.dbClient.Core
	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.CourseID != nil {
		query = query.Where("course_id = ?", options.CourseID)
	}
	if options.VideoID != nil {
		query = query.Where("video_id = ?", options.VideoID)
	}
	if options.ID != nil {
		query = query.Where("id = ?", options.ID)
	}
	tx := query.First(&question)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return question, nil
}
