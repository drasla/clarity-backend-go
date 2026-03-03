package service

import (
	"context"
	"errors"
	"strings"
	"tower/graph/model"
	"tower/model/maindb"
	"tower/pkg/fnError"
	"tower/repository"

	"gorm.io/gorm"
)

type EmailTemplateService interface {
	Create(ctx context.Context, input model.CreateEmailTemplateInput) (*maindb.EmailTemplate, error)
	FindById(ctx context.Context, id int) (*maindb.EmailTemplate, error)
	FindMany(ctx context.Context, page model.PageInput, search *model.EmailTemplateSearchInput) ([]maindb.EmailTemplate, int64, error)
	Update(ctx context.Context, id int, input model.ModifyEmailTemplateInput) (*maindb.EmailTemplate, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type emailTemplateService struct {
	repo repository.EmailTemplateRepository
}

func NewEmailTemplateService(repo repository.EmailTemplateRepository) EmailTemplateService {
	return &emailTemplateService{repo: repo}
}

func (s *emailTemplateService) Create(ctx context.Context, input model.CreateEmailTemplateInput) (*maindb.EmailTemplate, error) {
	existing, err := s.repo.FindByCode(ctx, input.TemplateCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fnError.NewInternalError(err, "템플릿 코드 중복 검사 중 오류가 발생했습니다.")
	}
	if existing != nil {
		return nil, fnError.NewBadRequest("이미 사용 중인 템플릿 코드입니다. 다른 코드를 입력해주세요.")
	}

	template := &maindb.EmailTemplate{
		TemplateCode: input.TemplateCode,
		Subject:      input.Subject,
		HTML:         input.HTML,
		Design:       input.Design,
	}

	if input.Variables != nil {
		template.Variables = *input.Variables
	}
	if input.Description != nil {
		template.Description = *input.Description
	}

	if err := s.repo.Create(ctx, template); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fnError.NewConflict("이미 존재하는 템플릿 코드입니다.")
		}
		return nil, fnError.NewInternalError(err, "템플릿 생성 중 오류가 발생했습니다.")
	}

	return template, nil
}

func (s *emailTemplateService) FindById(ctx context.Context, id int) (*maindb.EmailTemplate, error) {
	template, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fnError.NewNotFound("요청하신 이메일 템플릿을 찾을 수 없습니다.")
		}
		return nil, fnError.NewInternalError(err, "이메일 템플릿 조회 중 오류가 발생했습니다.")
	}
	return template, nil
}

func (s *emailTemplateService) FindMany(ctx context.Context, page model.PageInput, search *model.EmailTemplateSearchInput) ([]maindb.EmailTemplate, int64, error) {
	return s.repo.FindAll(ctx, page.Page, page.Size, search)
}

func (s *emailTemplateService) Update(ctx context.Context, id int, input model.ModifyEmailTemplateInput) (*maindb.EmailTemplate, error) {
	template, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, fnError.NewNotFound("수정할 템플릿을 찾을 수 없습니다.")
	}

	// 헬퍼 함수로 부분 업데이트
	mapModifyInputToTemplate(template, input)

	if err := s.repo.Update(ctx, template); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fnError.NewConflict("이미 사용 중인 템플릿 코드입니다.")
		}
		return nil, fnError.NewInternalError(err, "템플릿 수정 중 오류가 발생했습니다.")
	}
	return template, nil
}

func (s *emailTemplateService) Delete(ctx context.Context, id int) (bool, error) {
	_, err := s.repo.FindById(ctx, id)
	if err != nil {
		return false, fnError.NewNotFound("삭제할 이메일 템플릿을 찾을 수 없습니다.")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return false, fnError.NewInternalError(err, "이메일 템플릿 삭제 중 오류가 발생했습니다.")
	}

	return true, nil
}

func mapModifyInputToTemplate(template *maindb.EmailTemplate, input model.ModifyEmailTemplateInput) {
	if input.TemplateCode != nil {
		template.TemplateCode = *input.TemplateCode
	}
	if input.Subject != nil {
		template.Subject = *input.Subject
	}
	if input.HTML != nil {
		template.HTML = *input.HTML
	}
	if input.Design != nil {
		template.Design = *input.Design
	}
	if input.Variables != nil {
		template.Variables = *input.Variables
	}
	if input.Description != nil {
		template.Description = *input.Description
	}
}
