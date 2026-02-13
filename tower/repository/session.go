package repository

import (
	"context"
	"errors"
	"time"
	"tower/model/maindb"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

type SessionRepository interface {
	Create(ctx context.Context, token *maindb.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*maindb.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uint) error
	DeleteExpired(ctx context.Context) error
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, token *maindb.RefreshToken) error {
	return r.db.WithContext(ctx).
		Create(token).
		Error
}

func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*maindb.RefreshToken, error) {
	var refreshToken maindb.RefreshToken
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", token).
		First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&maindb.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *sessionRepository) RevokeAllForUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&maindb.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&maindb.RefreshToken{}).Error
}
