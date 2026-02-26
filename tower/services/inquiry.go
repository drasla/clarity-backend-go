package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"tower/graph/model"
	"tower/model/maindb"
	"tower/pkg/fnCrypto"
	"tower/pkg/fnError"
	"tower/pkg/fnMailer"
	"tower/pkg/fnMiddleware"
	"tower/pkg/fnNotifier"
	"tower/repository"

	"gorm.io/gorm"
)

type InquiryService interface {
	Create(ctx context.Context, input model.CreateInquiryInput) (*maindb.Inquiry, error)
	FindOneById(ctx context.Context, id int, password *string) (*maindb.Inquiry, error)
	FindManyPublic(ctx context.Context, page model.PageInput, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error)
	FindManyMyInquiries(ctx context.Context, page model.PageInput) ([]maindb.Inquiry, int64, error)
	FindManyForAdmin(ctx context.Context, page model.PageInput, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error)
	Modify(ctx context.Context, id int, input model.ModifyInquiryInput, password *string) (*maindb.Inquiry, error)
	Delete(ctx context.Context, id int, password *string) error
	Answer(ctx context.Context, id int, input model.AnswerInquiryInput) (*maindb.Inquiry, error)
}

type TelegramOptions struct {
	TelegramBotToken string
	TelegramChatID   string
}

type inquiryService struct {
	repo         repository.InquiryRepository
	templateRepo repository.EmailTemplateRepository
	opts         TelegramOptions
}

func NewInquiryService(repo repository.InquiryRepository, templateRepo repository.EmailTemplateRepository, opts TelegramOptions) InquiryService {
	return &inquiryService{
		repo:         repo,
		templateRepo: templateRepo,
		opts:         opts,
	}
}

func (s *inquiryService) Create(ctx context.Context, input model.CreateInquiryInput) (*maindb.Inquiry, error) {
	inquiry := &maindb.Inquiry{
		Category:    maindb.InquiryCategory(input.Category),
		Domain:      input.Domain,
		Title:       input.Title,
		Content:     input.Content,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Status:      maindb.InquiryStatusPending,
	}

	userID, ok := ctx.Value(fnMiddleware.UserIDKey).(uint)
	if ok && userID != 0 {
		inquiry.UserID = &userID
	} else {
		if input.NonMemberPw == nil || *input.NonMemberPw == "" {
			return nil, fnError.NewBadRequest("비회원 문의 시 비밀번호가 필요합니다.")
		}
		hashedPw, _ := fnCrypto.HashPassword(*input.NonMemberPw)
		inquiry.NonMemberPw = &hashedPw
	}

	if input.Attachments != nil {
		for _, fileInput := range input.Attachments {
			inquiry.Attachments = append(inquiry.Attachments, maindb.File{
				OriginalName: fileInput.OriginalName,
				StoredName:   fileInput.StoredName,
				URL:          fileInput.URL,
				Size:         int64(fileInput.Size),
				Extension:    fileInput.Extension,
				TargetType:   "inquiry",
			})
		}
	}

	if err := s.repo.Create(ctx, inquiry); err != nil {
		return nil, err
	}

	go func(inq *maindb.Inquiry) {
		token := s.opts.TelegramBotToken
		chatID := s.opts.TelegramChatID

		if token != "" && chatID != "" {
			msg := fmt.Sprintf(
				"🚨 <b>[새로운 1:1 문의 등록]</b>\n\n"+
					"<b>분류:</b> %s\n"+
					"<b>제목:</b> %s\n"+
					"<b>연락처:</b> %s\n"+
					"<b>이메일:</b> %s\n\n"+
					"관리자 페이지에서 확인해주세요!",
				inq.Category, inq.Title, inq.PhoneNumber, inq.Email,
			)

			fnNotifier.SendTelegramMessage(token, chatID, msg)
		}
	}(inquiry)

	return inquiry, nil
}

func (s *inquiryService) FindOneById(ctx context.Context, id int, password *string) (*maindb.Inquiry, error) {
	inquiry, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fnError.NewNotFound("문의글을 찾을 수 없습니다.")
		}
		return nil, fnError.NewInternalError(err, "조회 중 오류가 발생했습니다.")
	}

	role, _ := ctx.Value(fnMiddleware.RoleKey).(model.UserRole)
	if role == model.UserRoleAdmin {
		return inquiry, nil
	}

	if inquiry.UserID != nil {
		userID, _ := ctx.Value(fnMiddleware.UserIDKey).(uint)
		if *inquiry.UserID != userID {
			return nil, fnError.NewForbidden("본인의 문의글만 조회할 수 있습니다.")
		}
	} else {
		if password == nil || !fnCrypto.CheckPassword(*password, *inquiry.NonMemberPw) {
			return nil, fnError.NewForbidden("비밀번호가 일치하지 않습니다.")
		}
	}

	return inquiry, nil
}

func (s *inquiryService) FindManyPublic(ctx context.Context, page model.PageInput, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error) {
	return s.repo.FindPublicInquiries(ctx, page.Page, page.Size, search)
}

func (s *inquiryService) FindManyMyInquiries(ctx context.Context, page model.PageInput) ([]maindb.Inquiry, int64, error) {
	userID, ok := ctx.Value(fnMiddleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return nil, 0, fnError.NewUnauthorized("로그인이 필요합니다.")
	}
	return s.repo.FindMyInquiries(ctx, userID, page.Page, page.Size)
}

func (s *inquiryService) FindManyForAdmin(ctx context.Context, page model.PageInput, search *model.InquirySearchInput) ([]maindb.Inquiry, int64, error) {
	return s.repo.FindAll(ctx, page.Page, page.Size, search)
}

func (s *inquiryService) Modify(ctx context.Context, id int, input model.ModifyInquiryInput, password *string) (*maindb.Inquiry, error) {
	inquiry, err := s.FindOneById(ctx, id, password)
	if err != nil {
		return nil, err
	}

	if inquiry.Status == maindb.InquiryStatusCompleted {
		return nil, fnError.NewBadRequest("이미 답변이 완료된 문의는 수정할 수 없습니다.")
	}

	if input.Category != nil {
		inquiry.Category = maindb.InquiryCategory(*input.Category)
	}
	if input.Domain != nil {
		inquiry.Domain = input.Domain
	}
	if input.Title != nil {
		inquiry.Title = *input.Title
	}
	if input.Content != nil {
		inquiry.Content = *input.Content
	}
	if input.Email != nil {
		inquiry.Email = *input.Email
	}
	if input.PhoneNumber != nil {
		inquiry.PhoneNumber = *input.PhoneNumber
	}

	if input.Attachments != nil {
		var newFiles []maindb.File
		for _, f := range input.Attachments {
			newFiles = append(newFiles, maindb.File{
				OriginalName: f.OriginalName,
				StoredName:   f.StoredName,
				URL:          f.URL,
				Size:         int64(f.Size),
				Extension:    f.Extension,
				TargetType:   "inquiry",
			})
		}
		inquiry.Attachments = newFiles
	}

	if err := s.repo.Update(ctx, inquiry); err != nil {
		return nil, fnError.NewInternalError(err, "수정 중 오류가 발생했습니다.")
	}
	return inquiry, nil
}

func (s *inquiryService) Delete(ctx context.Context, id int, password *string) error {
	_, err := s.FindOneById(ctx, id, password)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *inquiryService) Answer(ctx context.Context, id int, input model.AnswerInquiryInput) (*maindb.Inquiry, error) {
	inquiry, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, fnError.NewNotFound("문의글을 찾을 수 없습니다.")
	}

	inquiry.Answer = &input.Answer
	inquiry.Status = maindb.InquiryStatus(input.Status)
	inquiry.AnsweredAt = new(time.Now())

	if err = s.repo.Update(ctx, inquiry); err != nil {
		return nil, fnError.NewInternalError(err, "답변 등록 중 오류가 발생했습니다.")
	}

	go func(inq *maindb.Inquiry, ans string) {
		if inq.Email == "" {
			return
		}

		bgCtx := context.Background()

		tmpl, err := s.templateRepo.FindByCode(bgCtx, "INQUIRY_REPLY")
		if err != nil {
			log.Printf("[이메일 발송 실패] 템플릿(INQUIRY_REPLY) 조회 오류: %v\n", err)
			return
		}

		templateContext := map[string]any{
			"Inquiry": inq,
			"Answer":  ans,
		}

		subject, _ := fnMailer.CompileTemplate(tmpl.Subject, templateContext)
		htmlBody, _ := fnMailer.CompileTemplate(tmpl.HTMLBody, templateContext)

		_ = fnMailer.Send("INQUIRY_REPLY", inq.Email, subject, htmlBody)
	}(inquiry, input.Answer)

	return inquiry, nil
}
