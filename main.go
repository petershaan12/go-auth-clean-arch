// @title Go Auth Clean Arch API
// @version 1.0
// @description API for You By Peter Shaan with PASETO authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and PASETO token.

package main

import "github.com/petershaan12/go-auth-clean-arch/internal/cmd"

func main() {
	cmd.Execute()
}
