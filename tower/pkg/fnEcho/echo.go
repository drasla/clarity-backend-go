package fnEcho

import (
	"errors"
	"log"
	"net/http"
	"tower/pkg/database"
	"tower/pkg/fnEnv"
	"tower/pkg/fnMiddleware"
	"tower/pkg/handler"
	"tower/repository"
	service "tower/services"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func StartEchoServer(db *database.Container, errHandler *handler.ErrorHandler) *http.Server {
	e := createEchoServer(db, errHandler)
	port := fnEnv.GetString("PORT", "8080")
	addr := ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}
	go func() {
		log.Printf("🚀 Starting Echo Server on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ HTTP Server Error: %v", err)
		}
	}()

	return srv
}

func createEchoServer(db *database.Container, errHandler *handler.ErrorHandler) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		ctx := c.Request().Context()
		appErr := errHandler.Handle(ctx, err)

		if appErr.Code == 0 {
			appErr.Code = http.StatusInternalServerError
		}

		response := map[string]interface{}{
			"success": false,
			"message": appErr.UserMessage,
			"code":    appErr.Code,
		}

		if c.Request().Method == http.MethodHead {
			if err := c.NoContent(appErr.Code); err != nil {
				e.Logger.Error("", err)
			}
			return
		}

		if err := c.JSON(appErr.Code, response); err != nil {
			e.Logger.Error("Failed to send error response: ", err)
		}
	}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2.0)))

	e.Use(fnMiddleware.JwtMiddleware())

	userRepo := repository.NewUserRepository(db.MainDB)
	sessionRepo := repository.NewSessionRepository(db.MainDB)
	verificationRepo := repository.NewVerificationRepository(db.MainDB)

	verificationService := service.NewVerificationService(verificationRepo)
	authService := service.NewAuthService(userRepo, sessionRepo, verificationService)
	userService := service.NewUserService(userRepo)
	fileService, err := service.NewS3Service()
	if err != nil {
		log.Fatalf("❌ Failed to initialize FileService: %v", err)
	}

	graphqlHandler := NewGraphQLServer(errHandler, authService, verificationService, userService, fileService)

	e.GET("/playground", echo.WrapHandler(playground.Handler("GraphQL", "/graphql")))
	e.POST("/graphql", echo.WrapHandler(graphqlHandler))
	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Tower is Healthy! (Echo v5)")
	})

	return e
}
