package fnMapper

import (
	"fmt"
	"tower/graph/model"
	"tower/model/maindb"
)

func UserToGraphQL(u *maindb.User) *model.User {
	if u == nil {
		return nil
	}
	res := &model.User{
		ID:             fmt.Sprint(u.ID),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Username:       u.Username,
		Email:          u.Email,
		Name:           u.Name,
		PhoneNumber:    u.PhoneNumber,
		LandlineNumber: u.LandlineNumber,
		Role:           model.UserRole(u.Role),
		Type:           model.UserType(u.Type),
		Status:         model.UserStatus(u.Status),
		AgreeEmail:     u.AgreeEmail,
		AgreeSms:       u.AgreeSMS,
	}
	if u.Type == maindb.TypeBusiness && u.BizRegNumber != nil {
		res.BizInfo = &model.BusinessInfo{
			BizRegNumber:  safeDefer(u.BizRegNumber),
			BizCeo:        safeDefer(u.BizCEO),
			BizType:       safeDefer(u.BizType),
			BizItem:       safeDefer(u.BizItem),
			BizZipCode:    safeDefer(u.BizZipCode),
			BizAddress1:   safeDefer(u.BizAddress1),
			BizAddress2:   u.BizAddress2,
			BizLicenseURL: u.BizLicenseURL,
		}
	}
	return res
}

func UsersToGraphQL(users []maindb.User) []*model.User {
	res := make([]*model.User, len(users))
	for i, u := range users {
		res[i] = UserToGraphQL(&u)
	}
	return res
}

//func ToGRPCUser(u *maindb.User)
