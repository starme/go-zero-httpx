package validation

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

// Validate runs struct validation using the shared validator instance and context.
func Validate(ctx context.Context, v any) error {
	validatorInstance := NewValidator(nil)

	if err := validatorInstance.Struct(v); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var ve ValidateError
			for _, fieldError := range validationErrors {
				if validatorInstance.trans != nil {
					ve = ve.AddString(fieldError.Translate(validatorInstance.trans))
				} else {
					ve = ve.AddString(fieldError.Error())
				}
			}
			return ve
		}

		return err
	}

	return nil
}
