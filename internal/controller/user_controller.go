package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/constants"
	"github.com/petershaan12/go-auth-clean-arch/resource/model"
	"github.com/petershaan12/go-auth-clean-arch/resource/response"
	"gorm.io/gorm"
)

type UserController struct {
	service model.UserMethodService
	env     library.Env
}

func NewUserController(service model.UserMethodService, env library.Env) *UserController {
	return &UserController{
		service: service,
		env:     env,
	}
}

// @Summary List User
// @Description API to get paginated list of users with optional search
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Page size (default: 10, max: 100)" default(10)
// @Param search query string false "Search by name, email, role, or department"
// @Success 200 {object} model.JsonResponseTotal{data=[]model.User} "List of users with pagination"
// @Failure 400 {object} model.JsonResponsError "Invalid query parameters"
// @Failure 401 {object} model.JsonResponsError "Unauthorized - Invalid or missing token"
// @Failure 403 {object} model.JsonResponsError "Forbidden - Insufficient permissions"
// @Failure 500 {object} model.JsonResponsError "Internal server error - Database or system error"
// @Router /user [get]
// @Security BearerAuth
func (s *UserController) List(c echo.Context) error {
	result, err := s.service.List(c.Get("ctx").(context.Context))
	if err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, 500, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterfaceTotal(c, 200, result, "success", len(result))
}

// @Summary Update User
// @Description Update existing user information
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param body body model.UpdateUserRequest true "User update data"
// @Success 201 {object} model.JsonResponse{data=model.User} "User updated successfully"
// @Failure 400 {object} model.JsonResponsError "Invalid input data"
// @Failure 404 {object} model.JsonResponsError "User not found"
// @Failure 500 {object} model.JsonResponsError "Internal server error"
// @Router /user/{id} [put]
// @Security BearerAuth
func (s *UserController) Update(c echo.Context) error {
	var req *model.UpdateUserRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, err.Error(), constants.BadRequest)
	}

	if err := c.Validate(req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, library.GetValueBetween(err.Error(), "Error:", "tag"), constants.BadRequest)
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.ResponseInterfaceError(c, http.StatusBadRequest, "invalid id parameter", constants.BadRequest)
	}

	tx := c.Get(constants.DBTransaction).(*gorm.DB)

	result, err := s.service.Update(c, tx, id, req)
	if err != nil {

		return response.ResponseInterfaceError(c, 500, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, http.StatusCreated, result, "Store User")
}

// @Summary Delete User
// @Description Soft delete user (mark as deleted)
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID to delete" minimum(1)
// @Success 200 {object} model.JsonResponse{data=string} "User deleted successfully"
// @Failure 400 {object} model.JsonResponsError "Invalid user ID format"
// @Failure 404 {object} model.JsonResponsError "User not found"
// @Failure 500 {object} model.JsonResponsError "Internal server error"
// @Router /user/{id} [delete]
// @Security BearerAuth
func (s *UserController) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.ResponseInterfaceError(c, http.StatusBadRequest, "invalid id parameter", constants.BadRequest)
	}

	tx := c.Get(constants.DBTransaction).(*gorm.DB)

	err = s.service.Delete(c, tx, id)
	if err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, 500, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, http.StatusOK, nil, "Delete User")
}

// @Summary Get User by ID
// @Description Get detailed user information by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} model.JsonResponse{data=model.User} "User details retrieved successfully"
// @Failure 400 {object} model.JsonResponsError "Invalid user ID format"
// @Failure 404 {object} model.JsonResponsError "User not found"
// @Failure 500 {object} model.JsonResponsError "Internal server error"
// @Router /user/{id} [get]
// @Security BearerAuth
func (s *UserController) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.ResponseInterfaceError(c, http.StatusBadRequest, "invalid id parameter", constants.BadRequest)
	}

	result, err := s.service.GetByID(c, id)
	if err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, 500, err.Error(), constants.InternalServerError)
	}
	return response.ResponseInterface(c, http.StatusOK, result, "Get User by ID")
}

// @Summary Create User
// @Description API to create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param auth body model.CreateUserRequest true "Authentication request"
// @Success 200 {object} model.JsonResponse{data=model.User} "User registered successfully"
// @Failure 400 {object} model.JsonResponsError "Bad request"
// @Failure 401 {object} model.JsonResponsError "Unauthorized"
// @Failure 500 {object} model.JsonResponsError "Internal error"
// @Router /user [post]
func (s *UserController) Create(c echo.Context) error {
	var req *model.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, err.Error(), constants.BadRequest)
	}

	if err := c.Validate(req); err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, http.StatusBadRequest, library.GetValueBetween(err.Error(), "Error:", "tag"), constants.BadRequest)
	}

	tx := c.Get(constants.DBTransaction).(*gorm.DB)

	result, err := s.service.Create(tx, req)
	if err != nil {
		log.Printf("Error in List: %v", err)
		return response.ResponseInterfaceError(c, 500, err.Error(), constants.InternalServerError)
	}

	return response.ResponseInterface(c, http.StatusCreated, result, "Store User")
}
