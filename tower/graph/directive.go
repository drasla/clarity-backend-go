package graph

import (
	"context"
	"tower/graph/model"
	"tower/pkg/fnError"
	"tower/pkg/fnMiddleware"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	_, err := fnMiddleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "접근 권한이 없습니다. (로그인 필요)",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	rawRole := ctx.Value(fnMiddleware.RoleKey)
	roleStr, ok := rawRole.(string)
	if !ok {
		return nil, fnError.NewForbidden("권한 정보를 읽을 수 없습니다.")
	}
	if roleStr != string(model.UserRoleAdmin) {
		return nil, fnError.NewForbidden("관리자 권한이 필요합니다.")
	}

	return next(ctx)
}
