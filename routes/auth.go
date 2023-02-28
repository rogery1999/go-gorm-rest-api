package routes

import (
	"github.com/labstack/echo/v4"
	authHandlers "github.com/rogery1999/go-gorm-rest-api/handlers/auth"
)

func setupAuthRoutes(g *echo.Group) {
	authG := g.Group("/auth")

	authG.POST("/login", authHandlers.AuthLogin)
}
