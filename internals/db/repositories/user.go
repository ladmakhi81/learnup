package repositories

import (
	"github.com/ladmakhi81/learnup/internals/db/entities"
	"gorm.io/gorm"
)

type UserRepo interface {
	Repository[entities.User]
}

type UserRepoImpl struct {
	RepositoryImpl[entities.User]
}

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{
		RepositoryImpl[entities.User]{
			db: db,
		},
	}
}
