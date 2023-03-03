package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rogery1999/go-gorm-rest-api/types"
	"github.com/rogery1999/go-gorm-rest-api/validation"
)

func ValidateRequestBody(requestBody interface{}) error {
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
