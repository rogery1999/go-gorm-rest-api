package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

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

// TODO
func setupJWT(e *echo.Echo) {
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		Skipper: func(c echo.Context) bool {
			return true
		},
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
