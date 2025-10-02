package library

import (
	"log"
	"net/http"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/petershaan12/go-auth-clean-arch/resource/constants"
)

type RequestHandler struct {
	Echo     *echo.Echo
	EchoFile *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func dateValidation(fl validator.FieldLevel) bool {
	_, err := time.Parse(constants.LayoutDate, fl.Field().String())
	return err == nil
}

// Update fungsi NewRequestHandler dengan fitur tambahan
func NewRequestHandler() *RequestHandler {
	engine := echo.New()
	engineFile := echo.New()

	// CORS untuk main engine
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
		AllowCredentials: true,
	}))

	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())

	validator := validator.New()
	err := validator.RegisterValidation("date", dateValidation)
	if err != nil {
		log.Fatal("Failed to register date validation:", err.Error())
	}
	engine.Validator = &CustomValidator{validator: validator}

	engineFile.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	engineFile.Use(FileCors)

	return &RequestHandler{
		Echo:     engine,
		EchoFile: engineFile,
	}
}

// Custom CORS untuk file handling
func FileCors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		allowList := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:8080": true,
		}

		if origin := req.Header.Get("Origin"); allowList[origin] {
			res.Header().Add("Access-Control-Allow-Origin", origin)
			res.Header().Add("Access-Control-Allow-Methods", "GET")
			res.Header().Add("Access-Control-Allow-Headers", "*")

			if req.Method != "OPTIONS" {
				err := next(c)
				if err != nil {
					c.Error(err)
				}
			}
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid CORS Origin")
		}

		return nil
	}
}
