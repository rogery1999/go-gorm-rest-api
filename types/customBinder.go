package types

import "github.com/labstack/echo/v4"

type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	return
}
