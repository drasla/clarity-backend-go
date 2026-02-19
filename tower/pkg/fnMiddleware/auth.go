package fnMiddleware

import (
	"context"
	"net/http"
	"strings"
	"tower/pkg/fnEnv"
	"tower/pkg/fnJwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type contextKey string

const UserIDKey contextKey = "UserID"

func JwtMiddleware() echo.MiddlewareFunc {
	secret := fnEnv.GetString("JWT_SECRET", "secret_key_needs_to_be_changed")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return next(c)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header format")
			}
			tokenString := parts[1]

			token, err := jwt.ParseWithClaims(tokenString, &fnJwt.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			if claims, ok := token.Claims.(*fnJwt.JwtCustomClaims); ok && token.Valid {
				ctx := context.WithValue(c.Request().Context(), UserIDKey, claims.UserID)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	userID, ok := val.(uint)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "invalid user id type in context")
	}
	return userID, nil
}
