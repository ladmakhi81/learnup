package db

import (
	"github.com/ladmakhi81/learnup/internals/db/repositories"
	"github.com/ladmakhi81/learnup/types"
	"gorm.io/gorm"
)

type UnitOfWork interface {
	Begin() (UnitOfWorkTx, error)
}

type UnitOfWorkTx interface {
	Commit() error
	Rollback() error

	AnswerRepo() repositories.AnswerRepo
	CartRepo() repositories.CartRepo
	CategoryRepo() repositories.CategoryRepo
	CommentRepo() repositories.CommentRepo
	CourseRepo() repositories.CourseRepo
	LikeRepo() repositories.LikeRepo
	NotificationRepo() repositories.NotificationRepo
	OrderRepo() repositories.OrderRepo
	OrderItemRepo() repositories.OrderItemRepo
	PaymentRepo() repositories.PaymentRepo
	QuestionRepo() repositories.QuestionRepo
	TransactionRepo() repositories.TransactionRepo
	UserRepo() repositories.UserRepo
	VideoRepo() repositories.VideoRepo
}

type UnitOfWorkImpl struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		db: db,
	}
}

func (svc UnitOfWorkImpl) Begin() (UnitOfWorkTx, error) {
	tx := svc.db.Begin()
	if tx.Error != nil {
		return nil, types.NewServerError(
			"Error in begin transaction",
			"UnitOfWorkImpl.Begin",
			tx.Error,
		)
	}
	return NewUnitOfWorkTx(tx), nil
}

type UnitOfWorkTxImpl struct {
	tx *gorm.DB
}

func NewUnitOfWorkTx(tx *gorm.DB) *UnitOfWorkTxImpl {
	return &UnitOfWorkTxImpl{
		tx: tx,
	}
}

func (svc UnitOfWorkTxImpl) Commit() error {
	err := svc.tx.Commit().Error
	return types.NewServerError(
		"Error in commit changes",
		"UnitOfWorkTxImpl.Commit",
		err,
	)
}

func (svc UnitOfWorkTxImpl) Rollback() error {
	err := svc.tx.Rollback().Error
	return types.NewServerError(
		"Error in rollback",
		"UnitOfWorkTxImpl.Rollback",
		err,
	)
}

func (svc UnitOfWorkTxImpl) AnswerRepo() repositories.AnswerRepo {
	return repositories.NewAnswerRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) CartRepo() repositories.CartRepo {
	return repositories.NewCartRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) CategoryRepo() repositories.CategoryRepo {
	return repositories.NewCategoryRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) CommentRepo() repositories.CommentRepo {
	return repositories.NewCommentRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) CourseRepo() repositories.CourseRepo {
	return repositories.NewCourseRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) LikeRepo() repositories.LikeRepo {
	return repositories.NewLikeRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) NotificationRepo() repositories.NotificationRepo {
	return repositories.NewNotificationRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) OrderRepo() repositories.OrderRepo {
	return repositories.NewOrderRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) OrderItemRepo() repositories.OrderItemRepo {
	return repositories.NewOrderItemRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) PaymentRepo() repositories.PaymentRepo {
	return repositories.NewPaymentRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) QuestionRepo() repositories.QuestionRepo {
	return repositories.NewQuestionRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) TransactionRepo() repositories.TransactionRepo {
	return repositories.NewTransactionRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) UserRepo() repositories.UserRepo {
	return repositories.NewUserRepo(svc.tx)
}
func (svc UnitOfWorkTxImpl) VideoRepo() repositories.VideoRepo {
	return repositories.NewVideoRepo(svc.tx)
}
