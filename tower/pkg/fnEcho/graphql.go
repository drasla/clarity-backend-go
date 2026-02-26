package fnEcho

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"context"
	localHandler "tower/pkg/fnError"
)

func NewGraphQLServer(
	errHandler *localHandler.ErrorHandler,
	execSchema graphql.ExecutableSchema,
) *handler.Server {
	srv := handler.New(execSchema)

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
