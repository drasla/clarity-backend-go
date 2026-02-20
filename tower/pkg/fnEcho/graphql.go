package fnEcho

import (
	"tower/graph"
	"tower/pkg/fnMiddleware"
	service "tower/services"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"context"
	localHandler "tower/pkg/handler"
)

func NewGraphQLServer(
	errHandler *localHandler.ErrorHandler,
	authService service.AuthService,
	verService service.VerificationService,
	userService service.UserService,
	fileService service.FileService,
) *handler.Server {
	c := graph.Config{
		Resolvers: &graph.Resolver{
			AuthService:         authService,
			VerificationService: verService,
			UserService:         userService,
			FileService:         fileService,
		},
	}

	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
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

	srv := handler.New(graph.NewExecutableSchema(c))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: 32 << 20,
		MaxMemory:     32 << 20,
	})
	srv.AddTransport(transport.Options{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		appErr := errHandler.Handle(ctx, e)

		return &gqlerror.Error{
			Message: appErr.UserMessage,
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": appErr.Code,
			},
		}
	})

	return srv
}
