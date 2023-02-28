package middlewares

import "github.com/labstack/echo/v4"

func SetupMiddlewares(e *echo.Echo) {
	JWTMiddleware(e)
	CORSMiddleware(e)
}
