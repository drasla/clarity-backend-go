package service

import (
	"context"
	"errors"
	"strings"
	"tower/graph/model"
	"tower/model/maindb"
	"tower/pkg/fnCrypto"
	"tower/pkg/fnError"
	"tower/pkg/fnMiddleware"
	"tower/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (*maindb.User, error)
	FindManyUserForAdmin(ctx context.Context, page model.PageInput, search *model.UserSearchInput) ([]maindb.User, int64, error)
	FindOneUserForAdmin(ctx context.Context, id uint) (*maindb.User, error)
	Modify(ctx context.Context, input model.ModifyUserInput) (*maindb.User, error)
	CheckPassword(ctx context.Context, password string) (bool, error)
	ChangePassword(ctx context.Context, input model.ChangePasswordInput) (bool, error)
	ModifyUserForAdmin(ctx context.Context, id uint, input model.ModifyUserForAdminInput) (*maindb.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id uint) (*maindb.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Status == maindb.StatusWithdrawn {
		return nil, errors.New("user has withdrawn")
	}
	if user.Status == maindb.StatusSuspended {
		return nil, errors.New("user is suspended")
	}

	return user, nil
}

func (s *userService) FindManyUserForAdmin(ctx context.Context, page model.PageInput, search *model.UserSearchInput) ([]maindb.User, int64, error) {
	return s.repo.FindAll(ctx, page.Page, page.Size, search)
}

func (s *userService) FindOneUserForAdmin(ctx context.Context, id uint) (*maindb.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fnError.NewNotFound("사용자 정보를 찾을 수 없습니다.")
	}

	return user, nil
}

func (s *userService) CheckPassword(ctx context.Context, password string) (bool, error) {
	userID, ok := ctx.Value(fnMiddleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return false, fnError.NewUnauthorized("로그인이 필요합니다.")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return false, fnError.NewNotFound("사용자 정보를 찾을 수 없습니다.")
	}

	isMatch := fnCrypto.CheckPassword(password, user.Password)
	return isMatch, nil
}

func (s *userService) Modify(ctx context.Context, input model.ModifyUserInput) (*maindb.User, error) {
	userID, ok := ctx.Value(fnMiddleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return nil, fnError.NewUnauthorized("로그인이 필요하거나 토큰이 만료되었습니다")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fnError.NewNotFound("사용자 정보를 찾을 수 없습니다.")
	}

	mapModifyInputToUser(user, input)
	if err := s.repo.Update(ctx, user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fnError.NewConflict("이미 사용 중인 이메일입니다.")
		}
		return nil, fnError.NewInternalError(err, "회원 정보 수정 중 오류가 발생했습니다.")
	}

	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (bool, error) {
	userID, ok := ctx.Value(fnMiddleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return false, fnError.NewUnauthorized("로그인이 필요합니다.")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return false, fnError.NewNotFound("사용자 정보를 찾을 수 없습니다.")
	}

	if !fnCrypto.CheckPassword(input.OldPassword, user.Password) {
		return false, fnError.NewBadRequest("기존 비밀번호가 일치하지 않습니다.")
	}

	if input.OldPassword == input.NewPassword {
		return false, fnError.NewBadRequest("새 비밀번호는 기존 비밀번호와 다르게 설정해야 합니다.")
	}

	hashedPw, err := fnCrypto.HashPassword(input.NewPassword)
	if err != nil {
		return false, fnError.NewInternalError(err, "비밀번호 암호화 중 오류가 발생했습니다.")
	}

	user.Password = hashedPw
	if err := s.repo.Update(ctx, user); err != nil {
		return false, fnError.NewInternalError(err, "비밀번호 변경 중 오류가 발생했습니다.")
	}

	return true, nil
}

func (s *userService) ModifyUserForAdmin(ctx context.Context, id uint, input model.ModifyUserForAdminInput) (*maindb.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fnError.NewNotFound("수정하려는 사용자를 찾을 수 없습니다.")
	}

	if err := mapModifyInputToUserForAdmin(user, input); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fnError.NewConflict("이미 사용 중인 이메일 또는 아이디입니다.")
		}
		return nil, fnError.NewInternalError(err, "관리자 권한 수정 중 서버 오류가 발생했습니다.")
	}

	return user, nil
}

func mapModifyInputToUser(user *maindb.User, input model.ModifyUserInput) {
	if input.Type != nil {
		newType := maindb.UserType(*input.Type)
		if user.Type == maindb.TypeBusiness && newType == maindb.TypePersonal {
			clearBizInfo(user)
		}
		user.Type = newType
	}

	mapCommonFields(user, input.Name, input.Email, input.PhoneNumber, input.LandlineNumber, input.AgreeEmail, input.AgreeSms)
	mapBizFields(user, input.BizRegNumber, input.BizCeo, input.BizType, input.BizItem, input.BizZipCode, input.BizAddress1, input.BizAddress2, input.BizLicenseURL)
}

func mapModifyInputToUserForAdmin(user *maindb.User, input model.ModifyUserForAdminInput) error {
	if input.Role != nil {
		user.Role = maindb.UserRole(*input.Role)
	}
	if input.Type != nil {
		user.Type = maindb.UserType(*input.Type)
	}
	if input.Status != nil {
		user.Status = maindb.UserStatus(*input.Status)
	}

	if input.Password != nil && *input.Password != "" {
		hashed, err := fnCrypto.HashPassword(*input.Password)
		if err != nil {
			return fnError.NewInternalError(err, "비밀번호 암호화 실패")
		}
		user.Password = hashed
	}

	mapCommonFields(user, input.Name, input.Email, input.PhoneNumber, input.LandlineNumber, input.AgreeEmail, input.AgreeSms)
	mapBizFields(user, input.BizRegNumber, input.BizCeo, input.BizType, input.BizItem, input.BizZipCode, input.BizAddress1, input.BizAddress2, input.BizLicenseURL)

	return nil
}

func mapCommonFields(user *maindb.User, name, email, phone, land *string, agreeEmail, agreeSms *bool) {
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if phone != nil {
		user.PhoneNumber = *phone
	}
	if land != nil {
		user.LandlineNumber = land
	}
	if agreeEmail != nil {
		user.AgreeEmail = *agreeEmail
	}
	if agreeSms != nil {
		user.AgreeSMS = *agreeSms
	}
}

func mapBizFields(user *maindb.User, reg, ceo, bType, item, zip, addr1, addr2, url *string) {
	if reg != nil {
		user.BizRegNumber = reg
	}
	if ceo != nil {
		user.BizCEO = ceo
	}
	if bType != nil {
		user.BizType = bType
	}
	if item != nil {
		user.BizItem = item
	}
	if zip != nil {
		user.BizZipCode = zip
	}
	if addr1 != nil {
		user.BizAddress1 = addr1
	}
	if addr2 != nil {
		user.BizAddress2 = addr2
	}
	if url != nil {
		user.BizLicenseURL = url
	}
}

func clearBizInfo(user *maindb.User) {
	user.BizRegNumber = nil
	user.BizCEO = nil
	user.BizType = nil
	user.BizItem = nil
	user.BizZipCode = nil
	user.BizAddress1 = nil
	user.BizAddress2 = nil
	user.BizLicenseURL = nil
}
