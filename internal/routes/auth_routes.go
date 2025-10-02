package routes

import (
	"github.com/petershaan12/go-auth-clean-arch/internal/controller"
	"github.com/petershaan12/go-auth-clean-arch/internal/middleware"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
)

type AuthRoutes struct {
	handler          library.RequestHandler
	authController   *controller.AuthController
	middlewareDB     *middleware.DBMiddleware
	pasetoMiddleware *middleware.PasetoMiddleware
}

func (s *AuthRoutes) Setup() {
	api := s.handler.Echo.Group("/auth")
	api.POST("/login", s.authController.Login, s.middlewareDB.HandlerDB())

	protected := api.Group("", s.pasetoMiddleware.Authorize())
	protected.POST("/logout", s.authController.Logout, s.middlewareDB.HandlerDB())
	protected.POST("/refresh", s.authController.RefreshToken, s.middlewareDB.HandlerDB())
}

func NewAuthRoutes(
	handler library.RequestHandler,
	authController *controller.AuthController,
	middlewareDB *middleware.DBMiddleware,
	pasetoMiddleware *middleware.PasetoMiddleware,
) *AuthRoutes {
	return &AuthRoutes{
		handler:          handler,
		authController:   authController,
		middlewareDB:     middlewareDB,
		pasetoMiddleware: pasetoMiddleware,
	}
}
