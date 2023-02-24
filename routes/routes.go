package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/handlers"
)

func SetupRoutes(e *echo.Echo) {

	// * Home middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, path := range []string{"/", ""} {
				if c.Path() == path {
					return c.String(http.StatusOK, "Hello world")
				}
			}
			return next(c)
		}
	})

	handlers.SetupUsersRoutes(e)
}
