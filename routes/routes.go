package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/handlers"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	handlers.SetupUsersRoutes(e)
}
