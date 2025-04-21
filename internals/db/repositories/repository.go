package repositories

import (
	"errors"
	"gorm.io/gorm"
)

type GetPaginatedOptions struct {
	Offset     *int
	Limit      *int
	Order      *string
	Relations  []string
	Conditions map[string]any
}

type GetAllOptions struct {
	Order      *string
	Relations  []string
	Conditions map[string]any
}

type Repository[T any] interface {
	Create(entity *T) error
	BatchInsert(entities []*T) error
	Delete(entity *T) error
	BatchDelete(entities []*T) error
	GetByID(id uint) (*T, error)
	GetOne(condition map[string]any) (*T, error)
	GetAll(options GetAllOptions) ([]*T, error)
	GetPaginated(options GetPaginatedOptions) ([]*T, int, error)
	Exist(condition map[string]any) (bool, error)
	Update(entity *T) error
}

type RepositoryImpl[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *RepositoryImpl[T] {
	return &RepositoryImpl[T]{db: db}
}

func (repo RepositoryImpl[T]) Create(entity *T) error {
	return repo.db.Create(&entity).Error
}

func (repo RepositoryImpl[T]) BatchInsert(entities []*T) error {
	return repo.db.Create(&entities).Error
}

func (repo RepositoryImpl[T]) Delete(entity *T) error {
	return repo.db.Delete(entity).Error
}

func (repo RepositoryImpl[T]) BatchDelete(entities []*T) error {
	return repo.db.Delete(&entities).Error
}

func (repo RepositoryImpl[T]) GetOne(condition map[string]any) (*T, error) {
	var entity *T
	tx := repo.db.Where(condition).First(&entity)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return entity, nil
}

func (repo RepositoryImpl[T]) GetAll(options GetAllOptions) ([]*T, error) {
	var entities []*T
	query := repo.db
	if options.Conditions != nil {
		query = query.Where(options.Conditions)
	}
	if options.Relations != nil {
		for _, relation := range options.Relations {
			query = query.Preload(relation)
		}
	}
	order := "created_at desc"
	if options.Order != nil {
		order = *options.Order
	}
	tx := query.Order(order).Find(&entities)
	if tx.Error != nil {
		return nil, query.Error
	}
	return entities, nil
}

func (repo RepositoryImpl[T]) GetPaginated(options GetPaginatedOptions) ([]*T, int, error) {
	var entities []*T
	var model T
	var count int64
	query := repo.db

	if options.Conditions != nil {
		query = query.Where(options.Conditions)
	}

	if err := query.Model(&model).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if options.Relations != nil {
		for _, relation := range options.Relations {
			query = query.Preload(relation)
		}
	}

	if options.Offset != nil && options.Limit != nil {
		query = query.Offset((*options.Offset) * (*options.Limit))
	}

	order := "created_at desc"
	if options.Order != nil {
		order = *options.Order
	}

	if err := query.Order(order).Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(count), nil
}

func (repo RepositoryImpl[T]) Exist(condition map[string]any) (bool, error) {
	var count int64
	var model T
	tx := repo.db.Model(&model).Where(condition).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}
	return count > 0, nil
}

func (repo RepositoryImpl[T]) GetByID(id uint) (*T, error) {
	return repo.GetOne(map[string]any{"id": id})
}

func (repo RepositoryImpl[T]) Update(entity *T) error {
	return repo.db.Updates(entity).Error
}
