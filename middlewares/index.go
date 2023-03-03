package middlewares

import "github.com/labstack/echo/v4"

func SetupMiddlewares(e *echo.Echo) {
	CORSMiddleware(e)
	JWTMiddleware(e)
}
