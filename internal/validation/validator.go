package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/harry713j/minurly/internal/errors"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.NewBadRequestErr("invalid input structure")
		}

		// Collect all field errors
		var fieldErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}

		return errors.NewBadRequestErr(fmt.Sprintf("validation failed: %v", fieldErrors))
	}
	return nil
}
