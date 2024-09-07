package validationservice

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.uber.org/multierr"
	"log"
)

var validate = validator.New()

func Validate(v interface{}) error {
	err := validate.Struct(v)
	var resultError error
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {

			for _, fieldErr := range validationErrs {
				singleErr := fmt.Errorf("error: Field '%s' failed validation with tag '%s'", fieldErr.Field(), fieldErr.Tag())
				resultError = multierr.Append(resultError, singleErr)
			}
			log.Printf("[validation] Validation errors: %v", resultError)
			return resultError
		}
		return err
	}
	return nil
}
