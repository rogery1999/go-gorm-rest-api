package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/routes"
	"github.com/rogery1999/go-gorm-rest-api/validation"
)

func main() {
	e := echo.New()
	validation.CreateValidator()

	setupLogs(e)
	setupEnvironmentVariables(e)

	// setupJWT(e)
	setupCORS(e)

	e.Static("/resources", "./static")

	routes.SetupRoutes(e)

	e.HTTPErrorHandler = errorHandler

	e.Logger.Info(fmt.Sprintf("Server running on http://%v:%v", os.Getenv("HOST"), os.Getenv("PORT")))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
