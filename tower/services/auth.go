package service

import (
	"context"
	"errors"
	"time"
	"tower/model/maindb"
	"tower/pkg/fnCrypto"
	"tower/pkg/fnEnv"
	"tower/pkg/fnJwt"
	"tower/repository"
)

type RegisterInput struct {
	Username       string
	Password       string
	Name           string
	Email          string
	PhoneNumber    string
	LandlineNumber *string
	AgreeEmail     bool
	AgreeSMS       bool
	BizInfo        *BusinessInput
}

type BusinessInput struct {
	BizRegNumber  string
	BizCEO        string
	BizType       string
	BizItem       string
	BizZipCode    string
	BizAddress1   string
	BizAddress2   string
	BizLicenseURL string
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*maindb.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // returns access, refresh
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Withdraw(ctx context.Context, userID uint) error
}

type authService struct {
	userRepo            repository.UserRepository
	sessionRepo         repository.SessionRepository
	verificationService VerificationService
	jwtSecret           string
}

func NewAuthService(u repository.UserRepository, s repository.SessionRepository, v VerificationService) AuthService {
	return &authService{
		userRepo:            u,
		sessionRepo:         s,
		verificationService: v,
		jwtSecret:           fnEnv.GetString("JWT_SECRET", "secret_key_needs_to_be_changed"),
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*maindb.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("already exists email")
	}

	isVerified, err := s.verificationService.IsVerified(ctx, input.PhoneNumber, maindb.VerifySMS)
	if err != nil {
		return nil, err
	}

	if !isVerified {
		return nil, errors.New("휴대폰 인증이 완료되지 않았습니다")
	}

	hashedPw, err := fnCrypto.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &maindb.User{
		Username:    input.Username,
		Password:    hashedPw,
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		AgreeEmail:  input.AgreeEmail,
		AgreeSMS:    input.AgreeSMS,
		Role:        maindb.RoleUser,
		Type:        maindb.TypePersonal,
		Status:      maindb.StatusActive,
	}

	if input.BizInfo != nil {
		user.Type = maindb.TypeBusiness
		user.BizRegNumber = &input.BizInfo.BizRegNumber
		user.BizCEO = &input.BizInfo.BizCEO
		user.BizType = &input.BizInfo.BizType
		user.BizItem = &input.BizInfo.BizItem
		user.BizZipCode = &input.BizInfo.BizZipCode
		user.BizAddress1 = &input.BizInfo.BizAddress1
		user.BizAddress2 = &input.BizInfo.BizAddress2
		user.BizLicenseURL = &input.BizInfo.BizLicenseURL
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	if user.Status == maindb.StatusWithdrawn {
		return "", "", errors.New("user has withdrawn")
	}
	if user.Status == maindb.StatusSuspended {
		return "", "", errors.New("user is suspended")
	}

	if !fnCrypto.CheckPassword(password, user.Password) {
		return "", "", errors.New("incorrect password")
	}

	accessToken, err := fnJwt.GenerateAccessToken(user.ID, string(user.Role), s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := fnJwt.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = s.sessionRepo.Create(ctx, &maindb.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14), // 2주
		ClientIP:  "127.0.0.1",                         // TODO: Context에서 실제 IP 추출 필요
		UserAgent: "Unknown",                           // TODO: Context에서 UserAgent 추출 필요
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	storedToken, err := s.sessionRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}
	if storedToken == nil {
		return "", "", errors.New("invalid refresh token")
	}

	if storedToken.IsRevoked {
		// 보안 경고: 폐기된 토큰을 사용하려 함 (해킹 시도 가능성)
		// TODO: 해당 유저의 모든 토큰을 날려버리는게 안전함 (선택사항)
		_ = s.sessionRepo.RevokeAllForUser(ctx, storedToken.UserID)
		return "", "", errors.New("revoked token reused")
	}
	if storedToken.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("token expired")
	}

	if storedToken.User.Status != maindb.StatusActive {
		return "", "", errors.New("user is not active")
	}

	newAccessToken, _ := fnJwt.GenerateAccessToken(storedToken.UserID, string(storedToken.User.Role), s.jwtSecret)
	newRefreshToken, _ := fnJwt.GenerateRefreshToken()

	_ = s.sessionRepo.Revoke(ctx, refreshToken)

	_ = s.sessionRepo.Create(ctx, &maindb.RefreshToken{
		UserID:    storedToken.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14),
		ClientIP:  storedToken.ClientIP,
		UserAgent: storedToken.UserAgent,
	})

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) Withdraw(ctx context.Context, userID uint) error {
	if err := s.userRepo.Withdraw(ctx, userID); err != nil {
		return err
	}

	return s.sessionRepo.RevokeAllForUser(ctx, userID)
}
