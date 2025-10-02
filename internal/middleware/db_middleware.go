package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/constants"
)

type DBMiddleware struct {
	handler library.RequestHandler
	db      library.Database
	env     library.Env
}

func NewDatabaseTrx(handler library.RequestHandler, db library.Database, env library.Env) *DBMiddleware {
	return &DBMiddleware{
		handler: handler,
		db:      db,
		env:     env,
	}
}

func (m *DBMiddleware) HandlerDB() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Println("beginning database transaction")
			txHandle := m.db.DB.WithContext(c.Request().Context()).Begin()
			defer func() {
				if r := recover(); r != nil {
					txHandle.Rollback()
					log.Println("rollback database transaction")
					return
				}
			}()

			// Put tx into Request context so repositories can pick it from ctx
			// (Kamu perlu define constants.DBTransaction di package constants)
			ctx := context.WithValue(c.Request().Context(), constants.DBTransaction, txHandle)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			// keep compatibility with existing code
			c.Set("ctx", req.Context())
			c.Set(constants.DBTransaction, txHandle)

			if err := next(c); err != nil {
				log.Println("commit err : ", err.Error())
				return err
			}

			// commit transaction on success status
			if statusInList(c.Response().Status,
				[]int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
				log.Println("commit database transaction")
				if err := txHandle.Commit().Error; err != nil {
					log.Println("commit err : ", err.Error())
					return err
				}
			} else {
				log.Println("rolling back database transaction")
				txHandle.Rollback()
				return nil
			}
			return nil
		}
	}
}

func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}
