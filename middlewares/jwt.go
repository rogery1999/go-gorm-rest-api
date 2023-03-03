package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(e *echo.Echo) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// * SKIP validation
			p := strings.Split(c.Path(), "/")
			lp := strings.Join([]string{p[len(p)-2], p[len(p)-1]}, "/")
			if lp == "auth/login" {
				return next(c)
			}

			// * Validating jwt
			ah := c.Request().Header.Get("Authorization")
			jwtS := strings.Split(ah, "Bearer ")[1]

			token, err := jwt.Parse(jwtS, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			// TODO
			// token.Claims.Valid()

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			expiresAt := claims["ExpiresAt"]
			df, ok := expiresAt.(float64)
			if !ok {
				c.Error(errors.New("invalid ExpiredAt date in the jwt"))
			}

			if int64(df)-time.Now().Unix() < 0 {
				c.Logger().Debug(fmt.Sprintf("HERE: %v", time.Now().Add(25*time.Hour)))

				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("this token has expired"))
			}

			// TODO: Add the claim data into the context creating a new context type
			return next(c)
		}
	})
}
