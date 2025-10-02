package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/petershaan12/go-auth-clean-arch/internal/token"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/constants"
	"github.com/petershaan12/go-auth-clean-arch/resource/model"
	"github.com/petershaan12/go-auth-clean-arch/resource/response"
)

type AuthController struct {
	service     model.AuthMethodService
	serviceUser model.UserMethodService
	env         library.Env
}

func NewAuthController(service model.AuthMethodService, serviceUser model.UserMethodService, env library.Env) *AuthController {
	return &AuthController{
		service:     service,
		serviceUser: serviceUser,
		env:         env,
	}
}

// @Summary Login
// @Description API to authenticate user and generate paseto token
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body model.AuthReq true "Authentication request"
// @Success 200 {object} model.JsonResponse{data=model.TokenOutput} "Authentication response with paseto token"
// @Failure 400 {object} model.JsonResponsError "Bad request"
// @Failure 401 {object} model.JsonResponsError "Unauthorized"
// @Failure 500 {object} model.JsonResponsError "Internal error"
// @Router /auth/login [post]
func (a *AuthController) Login(c echo.Context) error {
	var req *model.AuthReq
	if err := c.Bind(&req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, err.Error(), constants.BadRequest)
	}

	if err := c.Validate(req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, err.Error(), constants.BadRequest)
	}

	fmt.Println("User login attempt:", req.Email, "from IP:", req.IPAddress)

	auth, err := a.service.Login(c.Request().Context(), req)
	if err != nil {
		log.Printf("Error in List: %v", err)
		// to log who is requesting the token
		// even if authentication fails we still log the request
		return response.ResponseInterfaceError(c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
	}

	result, err := a.service.GenerateToken(c.Request().Context(), auth)
	if err != nil {
		log.Printf("Error in List: %v", err)
		c.Set("ip_address", req.IPAddress)
		return response.ResponseInterfaceError(c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, 200, result, "Auth")
}

// @Summary Refresh Token
// @Description API to refresh paseto token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body model.RefreshTokenReq true "Refresh token request"
// @Success 200 {object} model.JsonResponse{data=model.TokenOutput} "New access token"
// @Failure 400 {object} model.JsonResponsError "Bad request"
// @Failure 401 {object} model.JsonResponsError "Invalid refresh token"
// @Failure 500 {object} model.JsonResponsError "Internal error"
// @Router /auth/refresh [post]
func (a *AuthController) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	var req model.RefreshTokenReq
	if err := c.Bind(&req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, err.Error(), constants.BadRequest)
	}

	payload, err := a.service.VerifyRefreshToken(ctx, &req)
	if err != nil {
		log.Printf("Error in RefreshToken: %v", err)
		return response.ResponseInterfaceError(c, http.StatusUnauthorized, "Invalid refresh token: "+err.Error(), constants.Unauthorized)
	}

	user, err := a.serviceUser.GetByID(c, payload.UserId)
	if err != nil {
		log.Printf("Error in GetUserById: %v", err)
		return response.ResponseInterfaceError(c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
	}

	tokenOutput, err := a.service.GenerateToken(ctx, user)
	if err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, 200, tokenOutput, "Token Refreshed")
}

// @Summary Logout
// @Description API to logout user and invalidate all user sessions
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} model.JsonResponse{data=string} "Logout successful"
// @Failure 401 {object} model.JsonResponsError "Unauthorized"
// @Failure 500 {object} model.JsonResponsError "Internal error"
// @Router /auth/logout [post]
// @Security BearerAuth
func (a *AuthController) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	payload := c.Get("data_paseto").(*token.Payload)

	if err := a.service.Logout(ctx, payload); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusInternalServerError, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, 200, "Logout successful", "Logout")
}
