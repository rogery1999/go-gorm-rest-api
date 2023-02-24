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
	"github.com/rogery1999/go-gorm-rest-api/utils"
	"github.com/rogery1999/go-gorm-rest-api/validation"
)

type CreateUserRequestBody struct {
	Name string `json:"name" validate:"required"`
	Age  uint   `json:"age" validate:"gte=0,lte=130"`
}

func findUserById(c echo.Context) error {
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

func getAllUsers(c echo.Context) error {
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

func createUser(c echo.Context) error {
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	data.UsersData = append(data.UsersData, *newUser)

	return c.NoContent(http.StatusCreated)
}

func updateUser(c echo.Context) error {
	requestBody := new(CreateUserRequestBody)
	if err := c.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body data, review your 'Content-Type' header")
	}

	err := validation.Validator.Struct(requestBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// TODO: Improve error response
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body data  ❌")
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

func deleteUser(c echo.Context) error {
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

func SetupUsersRoutes(c *echo.Group) {
	ug := c.Group("/users")

	// * Routes
	ug.GET("", getAllUsers)
	ug.GET("/:userId", findUserById)
	ug.POST("", createUser)
	ug.PATCH("/:userId", updateUser)
	ug.DELETE("/:userId", deleteUser)
}
