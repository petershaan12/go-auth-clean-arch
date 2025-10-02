/*
Copyright © 2025 Peter Shaan <petershaan12@gmail.com>
*/
package cmd

import (
	"log"

	_ "github.com/petershaan12/go-auth-clean-arch/docs"
	"github.com/petershaan12/go-auth-clean-arch/internal/controller"
	"github.com/petershaan12/go-auth-clean-arch/internal/middleware"
	"github.com/petershaan12/go-auth-clean-arch/internal/repository"
	"github.com/petershaan12/go-auth-clean-arch/internal/routes"
	"github.com/petershaan12/go-auth-clean-arch/internal/service"
	"github.com/petershaan12/go-auth-clean-arch/internal/token"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the authentication server here",
	Long:  "Start the authentication server here.",
	Run:   server,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func server(cmd *cobra.Command, args []string) {
	env := library.ModuleConfig()
	db, _ := library.GetDatabase()
	requestHandler := library.NewRequestHandler()

	dbMiddleware := middleware.NewDatabaseTrx(*requestHandler, db, env)

	userRepo := repository.NewUserRepository(db)

	tokenMaker, err := token.NewPaseto(env.Paseto.Key, userRepo)
	if err != nil {
		log.Fatal("cannot create token maker: ", err)
	}
	pasetoMiddleware := middleware.NewPasetoTrx(tokenMaker)

	// User
	userService := service.NewUserService(userRepo, env)
	userController := controller.NewUserController(userService, env)

	// Auth
	authService := service.NewAuthService(userRepo, env, tokenMaker) // atau repo khusus kalau ada
	authController := controller.NewAuthController(authService, userService, env)

	routes.SetupRoutes(*requestHandler, userController, authController, dbMiddleware, pasetoMiddleware)

	// Swagger UI route disabled until docs package is generated with `swag init`
	requestHandler.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		if err := requestHandler.EchoFile.Start(":" + env.FileServerPort); err != nil {
			log.Println("file server error:", err.Error())
		}
	}()

	if err := requestHandler.Echo.Start(":" + env.Port); err != nil {
		log.Println("echo start server error:", err.Error())
	}
}
