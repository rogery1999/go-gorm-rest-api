package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/handlers"
)

func helloWorldHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

func SetupRoutes(e *echo.Echo) {
	gV1 := e.Group("/api/v1")

	// * Home middleware
	gV1.GET("", helloWorldHandler)
	gV1.GET("/", helloWorldHandler)

	handlers.SetupUsersRoutes(gV1)
}
