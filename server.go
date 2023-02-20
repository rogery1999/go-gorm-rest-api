package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/routes"
)

func main() {
	e := echo.New()

	e.Static("/resources", "./static")
	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8100"))
}
