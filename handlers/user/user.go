package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rogery1999/go-gorm-rest-api/data"
	"github.com/rogery1999/go-gorm-rest-api/models"
	"github.com/rogery1999/go-gorm-rest-api/types"
	"github.com/rogery1999/go-gorm-rest-api/utils"
	"github.com/rogery1999/go-gorm-rest-api/validation"
)

func FindUserById(c echo.Context) error {
	userId := c.Param("userId")
	fmt.Println("userId", userId)

	for _, user := range data.UsersData {
		id, err := strconv.Atoi(userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid identifier")
		}

		if user.Id == uint64(id) {
			return c.JSON(http.StatusOK, user)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("User with id %s not found", userId))
}

func GetAllUsers(c echo.Context) error {
	c.Logger().Info("Get all users")
	if len(c.QueryParams()) == 0 {
		return c.JSON(http.StatusOK, data.UsersData)
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

	offsetItems := (page - 1) * quantity
	return c.JSON(http.StatusOK, utils.SplitSlice(data.UsersData, uint(offsetItems), uint(offsetItems+quantity)))
}

func CreateUser(c echo.Context) error {
	requestBody := new(createUserRequestBody)
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	c.Logger().Debug(requestBody)
	if err := validateRequestBody(requestBody); err != nil {
		c.Error(err)
		return nil
	}

	newUser := models.User{
		Id:   data.UsersData[len(data.UsersData)-1].Id + 1,
		Name: requestBody.Name,
		Age:  requestBody.Age,
	}

	data.UsersData = append(data.UsersData, newUser)

	return c.NoContent(http.StatusCreated)
}

func UpdateUser(c echo.Context) error {
	requestBody := new(updateUserRequestBody)
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body data, review your 'Content-Type' header")
	}

	if err := validateRequestBody(requestBody); err != nil {
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

// TODO: refactor this
func validateRequestBody(requestBody interface{}) error {
	err := validation.Validator.Struct(requestBody)
	if err != nil {
		errors := make(map[string]string)

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("on %s field expect an %s but receive %v", err.Field(), err.Tag(), err.Value())
			errors[err.StructField()] = errorMessage
		}

		return &types.CustomError{Status: http.StatusBadRequest, Body: map[string]interface{}{
			"errors": errors,
		}}
	}

	return nil
}
