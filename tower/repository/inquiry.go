package repository

import (
	"context"
	"tower/graph/model"
	"tower/model/maindb"

	"gorm.io/gorm"
)

type InquiryRepository interface {
	Create(ctx context.Context, inquiry *maindb.Inquiry) error
	FindById(ctx context.Context, id int) (*maindb.Inquiry, error)
	FindPublicInquiries(ctx context.Context, page, size int, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error)
	FindMyInquiries(ctx context.Context, userId uint, page, size int) ([]maindb.Inquiry, int64, error)
	FindAll(ctx context.Context, page, size int, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error)
	Update(ctx context.Context, inquiry *maindb.Inquiry) error
	Delete(ctx context.Context, id int) error
}

type inquiryRepository struct {
	db *gorm.DB
}

func NewInquiryRepository(db *gorm.DB) InquiryRepository {
	return &inquiryRepository{db: db}
}

func (r *inquiryRepository) Create(ctx context.Context, inquiry *maindb.Inquiry) error {
	return r.db.WithContext(ctx).Create(inquiry).Error
}

func (r *inquiryRepository) FindById(ctx context.Context, id int) (*maindb.Inquiry, error) {
	var inquiry maindb.Inquiry
	err := r.db.WithContext(ctx).Preload("Attachments").First(&inquiry, id).Error
	return &inquiry, err
}

func (r *inquiryRepository) FindPublicInquiries(ctx context.Context, page, size int, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error) {
	var inquiries []maindb.Inquiry
	var total int64

	query := r.db.WithContext(ctx).Model(&maindb.Inquiry{}).Where("user_id IS NULL")

	if search != nil {
		if search.Status != nil {
			query = query.Where("status = ?", *search.Status)
		}
		if search.Category != nil {
			query = query.Where("category = ?", *search.Category)
		}
		if search.Keyword != nil && *search.Keyword != "" {
			query = query.Where("title LIKE ? OR content LIKE ?", "%"+*search.Keyword+"%", "%"+*search.Keyword+"%")
		}
	}

	query.Count(&total)
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&inquiries).Error

	return inquiries, total, err
}

func (r *inquiryRepository) FindMyInquiries(ctx context.Context, userId uint, page, size int) ([]maindb.Inquiry, int64, error) {
	var inquiries []maindb.Inquiry
	var total int64

	query := r.db.WithContext(ctx).Model(&maindb.Inquiry{}).Where("user_id = ?", userId)
	query.Count(&total)

	offset := (page - 1) * size
	err := query.Preload("Attachments").Order("created_at DESC").Offset(offset).Limit(size).Find(&inquiries).Error
	return inquiries, total, err
}

func (r *inquiryRepository) FindAll(ctx context.Context, page, size int, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error) {
	var inquiries []maindb.Inquiry
	var total int64

	query := r.db.WithContext(ctx).Model(&maindb.Inquiry{})

	if search != nil {
		if search.Status != nil {
			query = query.Where("status = ?", *search.Status)
		}
		if search.Category != nil {
			query = query.Where("category = ?", *search.Category)
		}
		if search.Domain != nil && *search.Domain != "" {
			query = query.Where("domain LIKE ?", "%"+*search.Domain+"%")
		}
		if search.Keyword != nil && *search.Keyword != "" {
			query = query.Where("title LIKE ? OR content LIKE ?", "%"+*search.Keyword+"%", "%"+*search.Keyword+"%")
		}
	}

	query.Count(&total)
	offset := (page - 1) * size
	err := query.Preload("Attachments").Order("created_at DESC").Offset(offset).Limit(size).Find(&inquiries).Error

	return inquiries, total, err
}

func (r *inquiryRepository) Update(ctx context.Context, inquiry *maindb.Inquiry) error {
	return r.db.WithContext(ctx).Save(inquiry).Error
}

func (r *inquiryRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&maindb.Inquiry{}, id).Error
}
