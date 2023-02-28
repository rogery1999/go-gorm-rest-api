package routes

import (
	"github.com/labstack/echo/v4"
	userHandlers "github.com/rogery1999/go-gorm-rest-api/handlers/user"
)

func setupUsersRoutes(c *echo.Group) {
	userG := c.Group("/user")

	// * Routes
	userG.GET("", userHandlers.GetAllUsers)
	userG.GET("/:userId", userHandlers.FindUserById)
	userG.POST("", userHandlers.CreateUser)
	userG.PATCH("/:userId", userHandlers.UpdateUser)
	userG.DELETE("/:userId", userHandlers.DeleteUser)
}
