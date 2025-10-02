package routes

import (
	"github.com/petershaan12/go-auth-clean-arch/internal/controller"
	"github.com/petershaan12/go-auth-clean-arch/internal/middleware"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
)

type UserRoutes struct {
	handler          library.RequestHandler
	userController   *controller.UserController
	middlewareDB     *middleware.DBMiddleware
	pasetoMiddleware *middleware.PasetoMiddleware
}

func (s *UserRoutes) Setup() {
	api := s.handler.Echo.Group("/user")
	api.GET("", s.userController.List, s.middlewareDB.HandlerDB())
	protected := api.Group("", s.pasetoMiddleware.Authorize())
	protected.GET("", s.userController.List, s.middlewareDB.HandlerDB())
	api.POST("", s.userController.Create, s.middlewareDB.HandlerDB())
	protected.PATCH("/:id", s.userController.Update, s.middlewareDB.HandlerDB())
	protected.DELETE("/:id", s.userController.Delete, s.middlewareDB.HandlerDB())
	protected.GET("/:id", s.userController.GetByID, s.middlewareDB.HandlerDB())
}

func NewUserRoutes(
	handler library.RequestHandler,
	userController *controller.UserController,
	middlewareDB *middleware.DBMiddleware,
	pasetoMiddleware *middleware.PasetoMiddleware,
) *UserRoutes {
	return &UserRoutes{
		handler:          handler,
		userController:   userController,
		middlewareDB:     middlewareDB,
		pasetoMiddleware: pasetoMiddleware,
	}
}
