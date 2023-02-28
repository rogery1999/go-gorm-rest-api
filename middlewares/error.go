package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/types"
)

func ErrorMiddleware(err error, c echo.Context) {
	code := http.StatusInternalServerError

	if customErr, ok := err.(*types.CustomError); ok {
		code = int(customErr.Status)
	}

	c.Logger().Error(err)

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
	}

	c.JSON(code, map[string]string{
		"message": err.Error(),
	})
}
