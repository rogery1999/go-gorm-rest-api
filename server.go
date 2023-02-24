package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rogery1999/go-gorm-rest-api/routes"
)

func main() {
	e := echo.New()

	setupLogs(e)
	setupEnvironmentVariables(e)

	setupCORS(e)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// c.Logger().Debug("First middleware")
			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// c.Logger().Debug("Second middleware")
			return next(c)
		}
	})

	e.Static("/resources", "./static")

	routes.SetupRoutes(e)

	e.HTTPErrorHandler = errorHandler

	e.Logger.Info(fmt.Sprintf("Server running on http://%v:%v", os.Getenv("HOST"), os.Getenv("PORT")))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}

func setupEnvironmentVariables(e *echo.Echo) {
	for _, arg := range os.Args[1:] {
		envArg := strings.Split(arg, "=")
		if len(envArg) == 2 && envArg[0] == "ENV" {
			os.Setenv(envArg[0], envArg[1])
			break
		}
	}

	envFile := ""
	switch os.Getenv("ENV") {
	case "DEV":
		envFile = ".dev"
	case "PROD":
		envFile = ".prod"
	case "TEST":
		envFile = ".test"
	}

	content, err := os.ReadFile(fmt.Sprintf("./%v.env", envFile))
	if err != nil {
		e.Logger.Fatal("Could not read environment file", err)
	}

	for _, variable := range strings.Split(string(content), "\n") {
		envData := strings.Split(variable, "=")
		os.Setenv(envData[0], envData[1])
	}

	e.Logger.Info("Environment variables successfully added")
}

func setupLogs(e *echo.Echo) {
	// Header log setup
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	e.Logger.SetLevel(log.DEBUG)

	// e.Logger.Error("Error Logs Working!")
	// e.Logger.Info("Info Logs Working!")
	// e.Logger.Debug("Debug Logs Working!")
	// e.Logger.Fatal("Fatal Logs Working!")

	e.Logger.Info("Logs setup and running")
}

func setupCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodPut, http.MethodPatch},
	}))
}

func errorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	c.Logger().Error(err)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
	}

	c.JSON(code, map[string]string{
		"message": err.Error(),
	})
}
