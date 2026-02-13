package server

import (
	"errors"
	"log"
	"net/http"
	"tower/pkg/database"
	"tower/pkg/env"
	"tower/pkg/handler"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func StartEchoServer(db *database.Container, errHandler *handler.ErrorHandler) *echo.Echo {
	e := createEchoServer(db, errHandler)
	port := env.GetString("PORT", "8080")
	go func() {
		log.Printf("🚀 Starting Echo Server on :%s", port)
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Echo Server Error: %v", err)
		}
	}()

	return e
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

	graphqlHandler := NewGraphQLServer(db, errHandler)

	e.GET("/playground", echo.WrapHandler(playground.Handler("GraphQL", "/query")))
	e.POST("/query", echo.WrapHandler(graphqlHandler))
	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Tower is Healthy! (Echo v5)")
	})

	return e
}
