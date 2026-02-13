package repository

import (
	"context"
	"errors"
	"tower/model/maindb"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(ctx context.Context, user *maindb.User) error
	FindByID(ctx context.Context, id uint) (*maindb.User, error)
	FindByEmail(ctx context.Context, email string) (*maindb.User, error)
	FindByUsername(ctx context.Context, username string) (*maindb.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *maindb.User) error
	Withdraw(ctx context.Context, userID uint) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *maindb.User) error {
	return r.db.WithContext(ctx).
		Create(user).
		Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*maindb.User, error) {
	var user maindb.User
	if err := r.db.WithContext(ctx).
		First(&user, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*maindb.User, error) {
	var user maindb.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*maindb.User, error) {
	var user maindb.User
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&maindb.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Update(ctx context.Context, user *maindb.User) error {
	return r.db.WithContext(ctx).
		Save(user).
		Error
}

func (r *userRepository) Withdraw(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&maindb.User{}).
			Where("id = ?", userID).
			Update("status", maindb.StatusWithdrawn).Error; err != nil {
			return err
		}

		if err := tx.
			Where("id = ?", userID).
			Delete(&maindb.User{}).
			Error; err != nil {
			return err
		}

		return nil
	})
}
