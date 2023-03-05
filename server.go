package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rogery1999/go-gorm-rest-api/config"
	"github.com/rogery1999/go-gorm-rest-api/middlewares"
	"github.com/rogery1999/go-gorm-rest-api/types"
	"github.com/rogery1999/go-gorm-rest-api/validation"
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
		if variable == "" {
			continue
		}
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

func setupServer(e *echo.Echo) {
	e.Binder = &types.CustomBinder{}
	validation.CreateValidator()
	e.Static("/resources", "./static")

	setupLogs(e)
	setupEnvironmentVariables(e)
	middlewares.SetupMiddlewares(e)
	config.SetupDB(e)
}
