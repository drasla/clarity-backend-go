package fnMapper

import "tower/graph/model"

func TokenToGraphQL(accessToken, refreshToken string) *model.Token {
	return &model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
