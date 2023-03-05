package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBGorm *gorm.DB

func connectToDB(e *echo.Echo) {
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_USER_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIMEZONE"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DryRun: false,
		Logger: newLogger,
	})

	if err != nil {
		e.Logger.Fatal(err)
	}

	DBGorm = db
	e.Logger.Info("Trying to connect to the database was successfully")
}

func setupTables(e *echo.Echo) {
	DBGorm.AutoMigrate(&models.User{})
}

func SetupDB(e *echo.Echo) {
	connectToDB(e)
	setupTables(e)
}
