package repository

import (
	"context"
	"tower/graph/model"
	"tower/model/maindb"

	"gorm.io/gorm"
)

type EmailTemplateRepository interface {
	Create(ctx context.Context, template *maindb.EmailTemplate) error
	FindById(ctx context.Context, id int) (*maindb.EmailTemplate, error)
	FindByCode(ctx context.Context, code string) (*maindb.EmailTemplate, error)
	FindAll(ctx context.Context, page, size int, search *model.EmailTemplateSearchInput) ([]maindb.EmailTemplate, int64, error)
	Update(ctx context.Context, template *maindb.EmailTemplate) error
	Delete(ctx context.Context, id int) error
}

type emailTemplateRepository struct {
	db *gorm.DB
}

func NewEmailTemplateRepository(db *gorm.DB) EmailTemplateRepository {
	return &emailTemplateRepository{db: db}
}

func (r *emailTemplateRepository) Create(ctx context.Context, template *maindb.EmailTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *emailTemplateRepository) FindById(ctx context.Context, id int) (*maindb.EmailTemplate, error) {
	var tmpl maindb.EmailTemplate
	err := r.db.WithContext(ctx).First(&tmpl, id).Error
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *emailTemplateRepository) FindByCode(ctx context.Context, code string) (*maindb.EmailTemplate, error) {
	var tmpl maindb.EmailTemplate
	err := r.db.WithContext(ctx).Where("template_code = ?", code).First(&tmpl).Error
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *emailTemplateRepository) FindAll(ctx context.Context, page, size int, search *model.EmailTemplateSearchInput) ([]maindb.EmailTemplate, int64, error) {
	var templates []maindb.EmailTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&maindb.EmailTemplate{})

	if search != nil && search.Keyword != nil && *search.Keyword != "" {
		k := "%" + *search.Keyword + "%"
		query = query.Where("(template_code LIKE ? OR subject LIKE ? OR description LIKE ?)", k, k, k)
	}

	query.Count(&total)

	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&templates).Error

	return templates, total, err
}

func (r *emailTemplateRepository) Update(ctx context.Context, template *maindb.EmailTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *emailTemplateRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&maindb.EmailTemplate{}, id).Error
}
