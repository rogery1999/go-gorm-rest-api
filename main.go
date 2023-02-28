package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/middlewares"
	"github.com/rogery1999/go-gorm-rest-api/routes"
)

func main() {
	e := echo.New()

	setupServer(e)

	routes.SetupRoutes(e)

	e.HTTPErrorHandler = middlewares.ErrorMiddleware

	e.Logger.Info(fmt.Sprintf("Server running on http://%v:%v", os.Getenv("HOST"), os.Getenv("PORT")))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
