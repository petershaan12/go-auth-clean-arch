package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/petershaan12/go-auth-clean-arch/internal/token"
)

type PasetoMiddleware struct {
	tokenMaker token.Maker
}

func NewPasetoTrx(tokenMaker token.Maker) *PasetoMiddleware {
	return &PasetoMiddleware{
		tokenMaker: tokenMaker,
	}
}

func (p *PasetoMiddleware) Authorize() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}

			// Verify token
			ctx := c.Request().Context()
			payload, err := p.tokenMaker.VerifyToken(ctx, tokenString)
			if err != nil {
				if errors.Is(err, token.ErrTokenRevoked) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
				}
				if errors.Is(err, token.ErrExpiredToken) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
			}

			// Set payload to context
			c.Set("data_paseto", payload)
			return next(c)
		}
	}
}
