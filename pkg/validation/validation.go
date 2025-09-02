package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// ValidateStructDetailed returns detailed validation errors (for logging)
func ValidateStructDetailed(s any) []string {
	var errors []string

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' Failed Validation '%s'",
				err.Field(), err.Tag()))
		}
	}

	return errors
}
