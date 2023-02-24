package validation

import "github.com/go-playground/validator/v10"

var Validator *validator.Validate

func CreateValidator() {
	if Validator == nil {
		Validator = validator.New()
	}
}
