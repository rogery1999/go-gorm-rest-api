package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/utils"
)

func AuthLogin(c echo.Context) error {
	c.Logger().Info("AuthLogin endpoint executed")
	body := new(ReqBodyAuthLogin)

	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := utils.ValidateRequestBody(body); err != nil {
		c.Error(err)
		return nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  "rogery",
		"ExpiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
	})

	jwtS, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ResBodyAuthLogin{
		Username: "rogery",
		JWT:      jwtS,
	})
}
