package service

import (
	"context"
	"errors"
	"log"
	"time"
	"tower/model/maindb"
	"tower/repository"
)

type VerificationService interface {
	SendCode(ctx context.Context, target string, vType maindb.VerificationType) error
	VerifyCode(ctx context.Context, target string, vType maindb.VerificationType, code string) error
	IsVerified(ctx context.Context, target string, vType maindb.VerificationType) (bool, error)
}

type verificationService struct {
	repo repository.VerificationRepository
}

func NewVerificationService(repo repository.VerificationRepository) VerificationService {
	return &verificationService{repo: repo}
}

func (s *verificationService) SendCode(ctx context.Context, target string, vType maindb.VerificationType) error {
	// TODO: 6자리 랜덤 코드 생성 (실제로는 util 패키지 사용)
	code := "123456"

	ver := &maindb.Verification{
		Target:     target,
		Type:       vType,
		Code:       code,
		IsVerified: false,
		ExpiresAt:  time.Now().Add(3 * time.Minute),
	}

	if err := s.repo.Create(ctx, ver); err != nil {
		return err
	}

	switch vType {
	case maindb.VerifySMS:
		log.Printf("📱 [SMS 발송] To: %s, Code: %s", target, code)
		// TODO: smsClient.Send(...)
	case maindb.VerifyEmail:
		log.Printf("📧 [Email 발송] To: %s, Code: %s", target, code)
		// TODO: emailClient.Send(...)
	}

	return nil
}

func (s *verificationService) VerifyCode(ctx context.Context, target string, vType maindb.VerificationType, code string) error {
	ver, err := s.repo.FindValidCode(ctx, target, vType)
	if err != nil {
		return err
	}
	if ver == nil {
		return errors.New("인증 시간이 만료되었거나 잘못된 요청입니다")
	}

	if ver.Code != code {
		return errors.New("인증 코드가 일치하지 않습니다")
	}

	return s.repo.MarkAsVerified(ctx, ver.ID)
}

func (s *verificationService) IsVerified(ctx context.Context, target string, vType maindb.VerificationType) (bool, error) {
	return s.repo.IsVerified(ctx, target, vType, 30*time.Minute)
}
