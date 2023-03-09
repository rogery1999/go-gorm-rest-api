package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/config"
	"github.com/rogery1999/go-gorm-rest-api/data"
	"github.com/rogery1999/go-gorm-rest-api/models"
	"github.com/rogery1999/go-gorm-rest-api/schemas"
	"github.com/rogery1999/go-gorm-rest-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func FindUserById(c echo.Context) error {
	c.Logger().Info("Find user by id executed")
	userId := c.Param("userId")
	c.Logger().Info(fmt.Sprintf("%T %v", userId, userId))
	userUUID, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id value")
	}

	user := models.User{UUID: uint64(userUUID)}
	result := config.DBGorm.Select("first_name", "last_name", "middle_name", "email", "birthday").Find(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("No user found with the user id %v", userId))
	}

	return c.JSON(http.StatusOK, *user.MapUserToUserInfoResponse())
}

func GetAllUsers(c echo.Context) error {
	c.Logger().Info("Get all users executed")
	users := make([]models.User, 0)
	if len(c.QueryParams()) == 0 {
		config.DBGorm.Select("first_name", "last_name", "middle_name", "email", "birthday").Find(&users)

		usersR := make([]schemas.UserInfoResponse, 0)
		for _, user := range users {
			usersR = append(usersR, *user.MapUserToUserInfoResponse())
		}

		return c.JSON(http.StatusOK, usersR)
	}

	// * Pagination
	qty, pg := strings.Trim(c.QueryParam("qty"), " "), strings.Trim(c.QueryParam("page"), "")
	qtyLength, pgLength := len([]rune(qty)), len([]rune(pg))

	if pgLength > 0 && qtyLength == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "when the queryParam 'page' is sent you must also sent the 'qty' queryParam")
	}

	quantity, err := strconv.Atoi(qty)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s is not a valid value", qty))
	}

	if qtyLength > 0 && pgLength == 0 {
		return c.JSON(http.StatusOK, utils.SplitSlice(data.UsersData, 0, uint(quantity)))
	}

	page, err := strconv.Atoi(pg)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s is not a valid value", pg))
	}
	offset := (page - 1) * quantity

	// TODO: test this with more users
	config.DBGorm.Select("first_name", "last_name", "middle_name", "email", "birthday").Limit(qtyLength).Offset(offset).Find(&users)

	usersR := make([]schemas.UserInfoResponse, 0)
	for _, user := range users {
		usersR = append(usersR, *user.MapUserToUserInfoResponse())
	}

	return c.JSON(http.StatusOK, usersR)
}

func CreateUser(c echo.Context) error {
	requestBody := new(schemas.CreateUserRequestBody)
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := utils.ValidateRequestBody(requestBody); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// * Find is the email is already registered
	userDuplicatedEmail := models.User{
		Email: requestBody.Email,
	}
	result := config.DBGorm.First(&userDuplicatedEmail)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusBadRequest, "this email is already registered")
	}

	birthday, err := time.Parse("2006-01-02", requestBody.Birthday)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user := models.User{
		FirstName:  requestBody.FirstName,
		MiddleName: requestBody.MiddleName,
		LastName:   requestBody.LastName,
		Email:      requestBody.Email,
		Password:   string(passwordHash),
		Birthday:   birthday,
	}

	result = config.DBGorm.Create(&user)

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.NoContent(http.StatusCreated)
}

func UpdateUser(c echo.Context) error {
	requestBody := new(schemas.UpdateUserRequestBody)
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body data, review your 'Content-Type' header")
	}

	if err := utils.ValidateRequestBody(requestBody); err != nil {
		c.Error(err)
		return nil
	}

	for idx, user := range data.UsersData {
		id, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		if user.Id == uint64(id) {
			data.UsersData[idx].Name = requestBody.Name
			if requestBody.Age > 0 {
				data.UsersData[idx].Age = requestBody.Age
			}
			return c.NoContent(http.StatusOK)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No user found with id %s", c.Param("userId")))
}

func DeleteUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid userId %s, expected an unsigned integer", c.Param("userId")))
	}

	userIndex := -1

	for i, v := range data.UsersData {
		if v.Id == uint64(userId) {
			userIndex = i
			break
		}
	}

	if userIndex < 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No user was found with the id %d", userId))
	}

	data.UsersData = append(data.UsersData[:userIndex], data.UsersData[(userIndex+1):]...)

	return c.NoContent(http.StatusAccepted)
}
