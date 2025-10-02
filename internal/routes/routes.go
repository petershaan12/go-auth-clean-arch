package routes

import (
	"github.com/petershaan12/go-auth-clean-arch/internal/controller"
	"github.com/petershaan12/go-auth-clean-arch/internal/middleware"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
)

type Route interface {
	Setup()
}

type Routes []Route

func SetupRoutes(
	handler library.RequestHandler,
	userController *controller.UserController,
	authController *controller.AuthController,
	dbMiddleware *middleware.DBMiddleware,
	pasetoMiddleware *middleware.PasetoMiddleware,
) {
	// Buat slice berisi semua route module
	routes := Routes{
		NewUserRoutes(handler, userController, dbMiddleware, pasetoMiddleware),
		NewAuthRoutes(handler, authController, dbMiddleware, pasetoMiddleware),
		// Tambah module lain di sini jika ada
	}

	// Loop dan panggil Setup untuk setiap route
	for _, r := range routes {
		r.Setup()
	}
}
