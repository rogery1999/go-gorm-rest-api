package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/config"
	"github.com/rogery1999/go-gorm-rest-api/models"
	"github.com/rogery1999/go-gorm-rest-api/schemas"
	"github.com/rogery1999/go-gorm-rest-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthLogin(c echo.Context) error {
	c.Logger().Info("AuthLogin endpoint executed")
	body := new(schemas.ReqBodyAuthLogin)

	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := utils.ValidateRequestBody(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user := models.User{
		Email: body.Email,
	}
	result := config.DBGorm.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusBadRequest, "Login failed, please review your email and/or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Login failed, please review your email and/or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authId":    user.UUID,
		"ExpiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
	})

	jwtS, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, schemas.ResBodyAuthLogin{
		User: fmt.Sprintf("%v %v %v", user.FirstName, user.MiddleName, user.LastName),
		JWT:  jwtS,
	})
}

// TODO
func AuthRefreshToken(c echo.Context) error {
	c.Logger().Info("AuthRefreshToken endpoint executed")

	return nil
}
