package repository

import (
	"context"
	"errors"
	"time"
	"tower/model/maindb"

	"gorm.io/gorm"
)

type verificationRepository struct {
	db *gorm.DB
}

type VerificationRepository interface {
	Create(ctx context.Context, v *maindb.Verification) error
	FindValidCode(ctx context.Context, target string, vType maindb.VerificationType) (*maindb.Verification, error)
	MarkAsVerified(ctx context.Context, id uint) error
	DeleteExpired(ctx context.Context) error
	IsVerified(ctx context.Context, target string, vType maindb.VerificationType, within time.Duration) (bool, error)
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) Create(ctx context.Context, v *maindb.Verification) error {
	return r.db.WithContext(ctx).
		Create(v).
		Error
}

func (r *verificationRepository) FindValidCode(ctx context.Context, target string, vType maindb.VerificationType) (*maindb.Verification, error) {
	var v maindb.Verification

	err := r.db.WithContext(ctx).
		Where("target = ? AND type = ?", target, vType).
		Where("expires_at > ?", time.Now()).
		Where("is_verified = ?", false).
		Order("created_at desc").
		First(&v).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *verificationRepository) MarkAsVerified(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&maindb.Verification{}).
		Where("id = ?", id).
		Update("is_verified", true).Error
}

func (r *verificationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&maindb.Verification{}).Error
}

func (r *verificationRepository) IsVerified(ctx context.Context, target string, vType maindb.VerificationType, within time.Duration) (bool, error) {
	var count int64
	limitTime := time.Now().Add(-within)

	err := r.db.WithContext(ctx).
		Model(&maindb.Verification{}).
		Where("target = ? AND type = ?", target, vType).
		Where("is_verified = ?", true).
		Where("updated_at > ?", limitTime).
		Count(&count).Error

	return count > 0, err
}
